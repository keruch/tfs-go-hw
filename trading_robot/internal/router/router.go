package router

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/keruch/tfs-go-hw/trading_robot/internal/domain"
	"github.com/keruch/tfs-go-hw/trading_robot/pkg/log"
)

type Subscriber interface {
	SubscribePairs(pairs ...string) error
	UnsubscribePairs(pairs ...string) error
}

type PriceQuantitySetter interface {
	SetPriceMultiplier(m float64)
	SetTradingQuantity(q int)
}

type Router struct {
	*mux.Router
	subscriber Subscriber
	options    PriceQuantitySetter
	logger     *log.Logger
}

func NewRouter(subscriber Subscriber, options PriceQuantitySetter, logger *log.Logger) *Router {
	r := &Router{
		subscriber: subscriber,
		options:    options,
		logger:     logger,
		Router:     mux.NewRouter(),
	}
	r.Methods(http.MethodPost).PathPrefix(domain.SubscribePair).HandlerFunc(r.postSubscribe)
	r.Methods(http.MethodPost).PathPrefix(domain.UnsubscribePair).HandlerFunc(r.postUnsubscribe)
	r.Methods(http.MethodPost).PathPrefix(domain.SetQuantity).HandlerFunc(r.postQuantity)
	r.Methods(http.MethodPost).PathPrefix(domain.SetMultiplier).HandlerFunc(r.postMultiplier)

	return r
}

func (r *Router) postSubscribe(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	ticker := vars[domain.PairVar]
	err := r.subscriber.SubscribePairs(ticker)
	if err != nil {
		r.logger.Errorf("%s endpoint: %s", domain.SubscribePair, err)
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func (r *Router) postUnsubscribe(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	ticker := vars[domain.PairVar]
	err := r.subscriber.UnsubscribePairs(ticker)
	if err != nil {
		r.logger.Errorf("%s endpoint: %s", domain.UnsubscribePair, err)
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func (r *Router) postQuantity(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	quantity := vars[domain.ValueVal]
	quantityInt, err := strconv.Atoi(quantity)
	if err != nil {
		r.logger.Errorf("%s endpoint: %s", domain.SetQuantity, err)
		writer.WriteHeader(http.StatusBadRequest)
	}
	r.options.SetTradingQuantity(quantityInt)
}

func (r *Router) postMultiplier(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	multiplier := vars[domain.ValueVal]
	multiplierFloat, err := strconv.ParseFloat(multiplier, 32)
	if err != nil {
		r.logger.Errorf("%s endpoint: %s", domain.SetMultiplier, err)
		writer.WriteHeader(http.StatusBadRequest)
	}
	r.options.SetPriceMultiplier(multiplierFloat)
}
