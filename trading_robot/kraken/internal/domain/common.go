package domain

import "errors"

const (
	OrderEndpoint = "/orders"

	PUBLIC_KEY  = "OLlMf8FawHE0FC4L7wZmvpxDEZcz184C0DbG6WgX3+ZoLb4s9EmtI9VR"
	PRIVATE_KEY = "tm+El0L45IDVJLoy05xRWf8Bu7IcyRpsrXwYe5GPtuLlKzcxE3rluKNf6oggf5E0jEUmgtGUo8WXywshBSEToZxM"
)

const (
	KrakenScheme = "https"
	KrakenHost   = "demo-futures.kraken.com"
	KrakenPath   = "/derivatives"
)

type OperationEndpoint string

const (
	KrakenCreateOrder     OperationEndpoint = "/api/v3/sendorder"
	KrakenOpenOrders      OperationEndpoint = "/api/v3/openorders"
	KrakenEditOrder       OperationEndpoint = "/api/v3/editorder" // does not supported because of orderId filed
	KrakenCancelOrder     OperationEndpoint = "/api/v3/cancelorder"
)

var ErrOperationNotFound = errors.New("given operation not found")
