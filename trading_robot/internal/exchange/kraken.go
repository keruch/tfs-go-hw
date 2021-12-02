package exchange

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"sync"

	rhttp "github.com/hashicorp/go-retryablehttp"
	"github.com/keruch/tfs-go-hw/trading_robot/config"
	"github.com/keruch/tfs-go-hw/trading_robot/internal/domain"
	"github.com/keruch/tfs-go-hw/trading_robot/internal/exchange/kraken"
	"github.com/keruch/tfs-go-hw/trading_robot/pkg/log"
	"github.com/keruch/tfs-go-hw/trading_robot/pkg/utils"
)

type KrakenExchange struct {
	logger *log.Logger

	client *rhttp.Client
	conn   *utils.RetryableWSConn

	mu    sync.RWMutex
	pairs map[string]bool
}

func NewKrakenExchange(logger *log.Logger) (*KrakenExchange, error) {
	u := url.URL{
		Scheme: kraken.WsScheme,
		Host:   kraken.Host,
		Path:   kraken.WsPath,
	}

	rwsconn := &utils.RetryableWSConn{
		Url:           u,
		MaxRetries:    kraken.MaxRetries,
		RequestHeader: nil,
	}

	_, err := rwsconn.RetryableDial()
	if err != nil {
		return nil, err
	}

	client := rhttp.NewClient()
	client.Logger = logger

	return &KrakenExchange{
		logger: logger,
		client: client,
		conn:   rwsconn,
		pairs:  make(map[string]bool),
	}, nil
}

func (k *KrakenExchange) sendOrder(order domain.Order, operation kraken.OperationEndpoint) (*kraken.ReceiveOrder, error) {
	queryMap, err := kraken.QueryByOperation(order, operation)
	if err != nil {
		return nil, err
	}

	req, err := k.createOrderRequest(operation, queryMap)
	if err != nil {
		return nil, err
	}

	resp, err := k.client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ro *kraken.ReceiveOrder
	err = json.Unmarshal(data, &ro)
	if err != nil {
		return nil, err
	}

	k.logger.Trace(string(data))

	return ro, nil
}

func (k *KrakenExchange) createOrderRequest(operation kraken.OperationEndpoint, queryParams kraken.QueryParams) (*rhttp.Request, error) {
	u := &url.URL{
		Scheme:   kraken.Scheme,
		Host:     kraken.Host,
		Path:     kraken.Path + string(operation),
		RawQuery: kraken.WsQuery,
	}

	q := u.Query()
	for key, val := range queryParams {
		q.Add(key, val)
	}
	u.RawQuery = q.Encode()

	method, err := kraken.GetMethodByOperation(operation)
	if err != nil {
		return nil, err
	}

	req, err := rhttp.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	token, err := kraken.GenerateToken(config.GetPrivateKey(), string(operation), q.Encode())
	if err != nil {
		return nil, err
	}
	req.Header.Set(kraken.Authent, token)
	req.Header.Set(kraken.APIKey, config.GetPublicKey())

	return req, nil
}

func (k *KrakenExchange) GetPrices(ctx context.Context) <-chan domain.Price {
	out := make(chan domain.Price)

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				k.logger.Info("Get Prices done")
				return
			default:
				_, data, reconnected, err := k.conn.ReadMessage()
				if err != nil {
					k.logger.Error(err)
					continue
				}

				if reconnected {
					if err = k.updateConnection(); err != nil {
						k.logger.Panic(err)
					}
				}

				k.logger.Trace(string(data))

				price, ok := utils.ValidateDataIsPrice(data)
				if ok {
					out <- price
				}
			}
		}
	}()

	return out
}

func (k *KrakenExchange) CloseConnection() error {
	return k.conn.Close()
}

func (k *KrakenExchange) SubscribePairs(pairs ...string) error {
	k.mu.RLock()
	pairsCount := len(k.pairs)
	k.mu.RUnlock()
	if pairsCount >= 1 {
		// TODO: delete hardcoded error
		return errors.New("can't subscribe to more than one ticker feed")
	}
	k.mu.Lock()
	for _, pair := range pairs {
		k.pairs[pair] = true
	}
	k.mu.Unlock()
	return k.sendRequest(kraken.SubscribeEvent, kraken.FeedType, pairs...)
}

func (k *KrakenExchange) UnsubscribePairs(pairs ...string) error {
	k.mu.Lock()
	for _, pair := range pairs {
		delete(k.pairs, pair)
	}
	k.mu.Unlock()
	return k.sendRequest(kraken.UnsubscribeEvent, kraken.FeedType, pairs...)
}

type request struct {
	Event      kraken.Event `json:"event"`
	Feed       kraken.Feed  `json:"feed"`
	ProductIDs []string     `json:"product_ids"`
}

func newRequest(event kraken.Event, feed kraken.Feed) request {
	return request{
		Event: event,
		Feed:  feed,
	}
}

func (k *KrakenExchange) sendRequest(event kraken.Event, feed kraken.Feed, tickers ...string) error {
	r := newRequest(event, feed)
	r.ProductIDs = append(r.ProductIDs, tickers...)
	reconnected, err := k.conn.WriteJSON(r)
	if err != nil {
		return err
	}

	if reconnected {
		if err = k.updateConnection(); err != nil {
			return err
		}
	}

	return nil
}

func (k *KrakenExchange) updateConnection() error {
	pairs := make([]string, 0)
	k.mu.RLock()
	for pair, _ := range k.pairs {
		pairs = append(pairs, pair)
	}
	k.mu.RUnlock()
	k.mu.Lock()
	k.pairs = make(map[string]bool)
	k.mu.Unlock()
	return k.SubscribePairs(pairs...)
}

func (k *KrakenExchange) CreateOrder(order domain.Order) (domain.CreateOrderResponse, error) {
	req, err := k.sendOrder(order, kraken.CreateOrder)
	if err != nil {
		return domain.CreateOrderResponse{}, err
	}

	return domain.CreateOrderResponse{
		OrderType:    order.OrderType,
		Symbol:       order.Symbol,
		Side:         order.Side,
		Size:         order.Size,
		LimitPrice:   order.LimitPrice,
		Result:       req.Result,
		Status:       req.SendStatus.Status,
		OrderID:      req.SendStatus.OrderID,
		ReceivedTime: req.SendStatus.ReceivedTime,
	}, nil
}

// GetOrders TODO: implement
func (k *KrakenExchange) GetOrders(order domain.Order) (*kraken.ReceiveOrder, error) {
	return k.sendOrder(order, kraken.OpenOrders)
}

// DeleteOrder TODO: implement
func (k *KrakenExchange) DeleteOrder(order domain.Order) (*kraken.ReceiveOrder, error) {
	return k.sendOrder(order, kraken.CancelOrder)
}
