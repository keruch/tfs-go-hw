package exchange

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"

	rhttp "github.com/hashicorp/go-retryablehttp"
	"github.com/keruch/tfs-go-hw/trading_robot/internal/exchange/kraken"
	"github.com/keruch/tfs-go-hw/trading_robot/pkg/log"
	"github.com/keruch/tfs-go-hw/trading_robot/pkg/utils"
)

type krakenExchange struct {
	logger *log.Logger

	client *rhttp.Client
	conn   *utils.RetryableWSConn

	shutdown chan bool

	mu    sync.RWMutex
	pairs map[string]bool
}

func NewKrakenExchange(logger *log.Logger) (Exchange, error) {
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

	return &krakenExchange{
		logger:   logger,
		client:   rhttp.NewClient(),
		conn:     rwsconn,
		shutdown: make(chan bool),
		pairs:    make(map[string]bool),
	}, nil
}

func (k *krakenExchange) SendOrder(order *kraken.SendOrder, operation kraken.OperationEndpoint) (*kraken.ReceiveOrder, error) {
	queryMap, err := QueryByOperation(order, operation)
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

	return ro, nil
}

func (k *krakenExchange) createOrderRequest(operation kraken.OperationEndpoint, queryParams kraken.QueryParams) (*rhttp.Request, error) {
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

	method, err := getMethodByOperation(operation)
	if err != nil {
		return nil, err
	}

	req, err := rhttp.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	token, err := generateToken(os.Getenv("PRIVATE_KEY"), string(operation), q.Encode())
	if err != nil {
		return nil, err
	}
	req.Header.Set(kraken.Authent, token)
	req.Header.Set(kraken.APIKey, os.Getenv("PUBLIC_KEY"))

	return req, nil
}

func getMethodByOperation(operation kraken.OperationEndpoint) (string, error) {
	switch {
	case operation == kraken.OpenOrders:
		return http.MethodGet, nil
	case operation == kraken.CreateOrder || operation == kraken.EditOrder || operation == kraken.CancelOrder:
		return http.MethodPost, nil
	default:
		return "", kraken.ErrOperationNotFound
	}
}

func generateToken(privateKey, endpoint, postData string) (string, error) {
	// step1
	step1 := postData + endpoint

	// step2
	sha := sha256.New()
	sha.Write([]byte(step1))
	step2 := sha.Sum(nil)

	// step 3
	step3, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", err
	}

	// step 4
	h := hmac.New(sha512.New, step3)
	h.Write(step2)
	step4 := h.Sum(nil)

	// step 5
	step5 := base64.StdEncoding.EncodeToString(step4)

	return step5, nil
}

func (k *krakenExchange) Listen() <-chan []byte {
	out := make(chan []byte)

	go func() {
		defer close(out)
		for {
			select {
			case <-k.shutdown:
				return
			default:
				_, data, reconnected, err := k.conn.ReadMessage()
				if err != nil {
					// TODO: add recover
					k.logger.Panic(err)
				}

				if reconnected {
					if err = k.updateConnection(); err != nil {
						k.logger.Panic(err)
					}
				}

				// TODO: unmarshal only if message is valid (text template)
				//var candleData kraken.CandleData
				//err = json.Unmarshal(rawTicker, &candleData)
				//if err != nil {
				//	//panic(err)
				//}

				out <- data
			}
		}
	}()

	return out
}

func (k *krakenExchange) Shutdown() error {
	// TODO: add context functionality
	//if err := k.closeConnection(); err != nil {
	//	return err
	//}
	k.shutdown <- true

	return k.conn.Close()
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

func (k *krakenExchange) SubscribePairs(pairs ...string) error {
	// TODO: add customization for feed type
	k.mu.Lock()
	for _, pair := range pairs {
		k.pairs[pair] = true
	}
	k.mu.Unlock()
	return k.sendRequest(kraken.SubscribeEvent, kraken.TickerFeed, pairs...)
}

func (k *krakenExchange) UnsubscribePairs(pairs ...string) error {
	// TODO: add customization for feed type
	k.mu.Lock()
	for _, pair := range pairs {
		delete(k.pairs, pair)
	}
	k.mu.Unlock()
	return k.sendRequest(kraken.UnsubscribeEvent, kraken.TickerFeed, pairs...)
}

func (k *krakenExchange) sendRequest(event kraken.Event, feed kraken.Feed, tickers ...string) error {
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

func (k *krakenExchange) updateConnection() error {
	return k.connectionControl(k.SubscribePairs)
}

func (k *krakenExchange) closeConnection() error {
	return k.connectionControl(k.UnsubscribePairs)
}

func (k *krakenExchange) connectionControl(operation func (pairs ...string) error) error {
	pairs := make([]string, 0)
	k.mu.RLock()
	for pair, _ := range k.pairs {
		pairs = append(pairs, pair)
	}
	k.mu.RUnlock()
	return operation(pairs...)
}

func (k *krakenExchange) CreateOrder(order *kraken.SendOrder) (*kraken.ReceiveOrder, error) {
	return k.SendOrder(order, kraken.CreateOrder)
}

func (k *krakenExchange) GetOrders(order *kraken.SendOrder) (*kraken.ReceiveOrder, error) {
	return k.SendOrder(order, kraken.OpenOrders)
}

func (k *krakenExchange) DeleteOrder(order *kraken.SendOrder) (*kraken.ReceiveOrder, error) {
	return k.SendOrder(order, kraken.CancelOrder)
}

func QueryByOperation(order *kraken.SendOrder, operation kraken.OperationEndpoint) (kraken.QueryParams, error) {
	switch operation {
	case kraken.CreateOrder:
		return kraken.QueryParams{
			kraken.OrderType:  order.OrderType,
			kraken.Symbol:     order.Symbol,
			kraken.Side:       order.Side,
			kraken.Size:       order.Size,
			kraken.LimitPrice: order.LimitPrice,
		}, nil

	case kraken.OpenOrders:
		return kraken.QueryParams{}, nil

	case kraken.CancelOrder:
		return kraken.QueryParams{
			kraken.OrderID: order.OrderID,
		}, nil

	default:
		return nil, kraken.ErrOperationNotFound
	}

}
