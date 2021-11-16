package repository

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/keruch/tfs-go-hw/trading_robot/data_processor/internal/domain"
)

type Repository interface {
	GetTickers() <-chan domain.Ticker
}

const (
	hostAddr = ":8081"

	tickersEndpoint = "/tickers"
)

type RestAPIRepo struct {
}

func (rar *RestAPIRepo) Listen() {
	r := mux.NewRouter()

	r.Methods("POST").PathPrefix(tickersEndpoint).HandlerFunc(postTicker)

	srv := &http.Server{
		Handler: r,
		Addr: hostAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func postTicker(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	log.Println(string(data))
}
