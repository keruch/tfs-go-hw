package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/keruch/tfs-go-hw/trading_robot/data_receiver/internal/domain"
	"github.com/keruch/tfs-go-hw/trading_robot/data_receiver/internal/repository"
)

type Service interface {
	Subscribe(ctx context.Context, tickers []string) <-chan domain.Ticker
}

type ReceiverService struct {
	repo repository.Repository
}

func NewReceiverService(repo repository.Repository) Service {
	return &ReceiverService{
		repo: repo,
	}
}

func (rs *ReceiverService) Subscribe(ctx context.Context, tickers []string) <-chan domain.Ticker {
	in, err := rs.repo.GetTickersData(ctx, tickers)

	out := make(chan domain.Ticker)

	go func() {
		defer close(out)
		for data := range in {
			log.Println(string(data))
			var ticker *domain.Ticker
			err = json.Unmarshal(data, &ticker)
			if err != nil {
				panic(err)
			}
			if ticker == nil {
				// TODO: delete hardcode
				panic("ticker is nil")
			}

			out <- *ticker
		}
	}()

	return out
}
