package repository

import (
	"net/http"

	rhttp "github.com/hashicorp/go-retryablehttp"
)

type Repository interface {
	Request(r *rhttp.Request) (*http.Response, error)
}

type KrakenExchange struct {
	client *rhttp.Client
}

func NewKrakenExchange() Repository{
	return &KrakenExchange{
		client: rhttp.NewClient(),
	}
}

func (ke *KrakenExchange) Request(r *rhttp.Request) (*http.Response, error) {
	resp, err := ke.client.Do(r)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
