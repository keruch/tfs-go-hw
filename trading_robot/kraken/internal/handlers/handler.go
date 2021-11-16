package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/keruch/tfs-go-hw/trading_robot/kraken/internal/domain"
	"github.com/keruch/tfs-go-hw/trading_robot/kraken/internal/service"
)

type Handler interface {
	PostOrder(w http.ResponseWriter, r *http.Request)
	GetOrders(w http.ResponseWriter, r *http.Request)
	PutOrder(w http.ResponseWriter, r *http.Request)
	DeleteOrder(w http.ResponseWriter, r *http.Request)
}

type KrakenHandler struct {
	service service.Service
}

func NewKrakenHandler(service service.Service) Handler {
	return &KrakenHandler{
		service: service,
	}
}

func (kh *KrakenHandler) PostOrder(w http.ResponseWriter, r *http.Request) {
	kh.processRequest(w, r, domain.KrakenCreateOrder)
	w.WriteHeader(http.StatusCreated)
}

func (kh *KrakenHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	kh.processRequest(w, r, domain.KrakenOpenOrders)
	w.WriteHeader(http.StatusOK)
}

func (kh *KrakenHandler) PutOrder(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (kh *KrakenHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	kh.processRequest(w, r, domain.KrakenCancelOrder)
	w.WriteHeader(http.StatusResetContent)
}

func (kh *KrakenHandler) processRequest(w http.ResponseWriter, r *http.Request, operation domain.OperationEndpoint) {
	var order domain.SendOrder
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		// TODO: handle clearly
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	query, err := queryByOperation(operation, order)
	if err != nil {
		// TODO: handle clearly
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := kh.service.OrderRequest(operation, query)
	if err != nil {
		// TODO: handle clearly
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	d, err := json.MarshalIndent(response, "", " ")
	if err != nil {
		// TODO: handle clearly
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(d)
	if err != nil {
		// TODO: handle clearly
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func queryByOperation(operation domain.OperationEndpoint, order domain.SendOrder) (domain.QueryParams, error) {
	switch operation {
	case domain.KrakenCreateOrder:
		return domain.QueryParams{
			domain.OrderType: order.OrderType,
			domain.Symbol: order.Symbol,
			domain.Side: order.Side,
			domain.Size: order.Size,
			domain.LimitPrice: order.LimitPrice,
		}, nil
	case domain.KrakenOpenOrders:
		return domain.QueryParams{}, nil
	case domain.KrakenEditOrder:
		return domain.QueryParams{
			domain.OrderID: order.OrderID,
			domain.Size: order.Size,
			domain.LimitPrice: order.LimitPrice,
		}, nil
	case domain.KrakenCancelOrder:
		return domain.QueryParams{
			domain.OrderID: order.OrderID,
		}, nil
	default:
		return nil, domain.ErrOperationNotFound
	}

}
