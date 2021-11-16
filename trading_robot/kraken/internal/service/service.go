package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	rhttp "github.com/hashicorp/go-retryablehttp"
	"github.com/keruch/tfs-go-hw/trading_robot/kraken/internal/domain"
	"github.com/keruch/tfs-go-hw/trading_robot/kraken/internal/repository"
)

type Service interface {
	OrderRequest(operation domain.OperationEndpoint, queryParams domain.QueryParams) (domain.ReceiveOrder, error)
}

type KrakenService struct {
	apiKey string
	repo   repository.Repository
}

func NewKrakenService(repo repository.Repository, key string) Service {
	return &KrakenService{
		apiKey: key,
		repo:   repo,
	}
}

func (ks *KrakenService) OrderRequest(operation domain.OperationEndpoint, queryParams domain.QueryParams) (domain.ReceiveOrder, error) {
	req, err := ks.createOrderRequest(operation, queryParams)
	if err != nil {
		return domain.ReceiveOrder{}, err
	}

	resp, err := ks.repo.Request(req)
	if err != nil {
		return domain.ReceiveOrder{}, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.ReceiveOrder{}, err
	}
	defer resp.Body.Close()

	var ro domain.ReceiveOrder
	err = json.Unmarshal(data, &ro)
	if err != nil {
		return domain.ReceiveOrder{}, err
	}

	return ro, nil
}

func (ks *KrakenService) createOrderRequest(operation domain.OperationEndpoint, queryParams domain.QueryParams) (*rhttp.Request, error) {
	u := &url.URL{
		Scheme: domain.KrakenScheme,
		Host:   domain.KrakenHost,
		Path:   domain.KrakenPath + string(operation),
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

	token, err := generateToken(domain.PRIVATE_KEY, string(operation), q.Encode())
	if err != nil {
		return nil, err
	}
	req.Header.Set(domain.Authent, token)
	req.Header.Set(domain.APIKey, ks.apiKey)

	return req, nil
}

func getMethodByOperation(operation domain.OperationEndpoint) (string, error) {
	switch {
	case operation == domain.KrakenOpenOrders:
		return http.MethodGet, nil
	case operation == domain.KrakenCreateOrder || operation == domain.KrakenEditOrder || operation == domain.KrakenCancelOrder:
		return http.MethodPost, nil
	default:
		return "", domain.ErrOperationNotFound
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
