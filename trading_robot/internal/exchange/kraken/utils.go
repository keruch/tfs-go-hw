package kraken

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"github.com/keruch/tfs-go-hw/trading_robot/internal/domain"
)

func GetMethodByOperation(operation OperationEndpoint) (string, error) {
	switch {
	case operation == OpenOrders:
		return http.MethodGet, nil
	case operation == CreateOrder || operation == EditOrder || operation == CancelOrder:
		return http.MethodPost, nil
	default:
		return "", ErrOperationNotFound
	}
}

func GenerateToken(privateKey, endpoint, postData string) (string, error) {
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

func QueryByOperation(order domain.Order, operation OperationEndpoint) (QueryParams, error) {
	switch operation {
	case CreateOrder:
		return QueryParams{
			OrderType:     order.OrderType,
			Symbol:        order.Symbol,
			Side:          order.Side,
			Size:          strconv.Itoa(order.Size),
			LimitPrice:    fmt.Sprintf("%.1f", order.LimitPrice),
			StopPrice:     fmt.Sprintf("%.1f", order.StopPrice),
			TriggerSignal: order.TriggerSignal,
		}, nil

	case OpenOrders:
		return QueryParams{}, nil

	case CancelOrder:
		return QueryParams{
			OrderID: order.OrderID,
		}, nil

	default:
		return nil, ErrOperationNotFound
	}
}
