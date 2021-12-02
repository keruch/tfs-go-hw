package repository

import (
	"context"

	"github.com/keruch/tfs-go-hw/hw3/internal/domain"
	"github.com/keruch/tfs-go-hw/hw3/internal/domain/generator"
)

type PricesData interface {
	GetPrices(ctx context.Context) <-chan domain.CandleData
}

type generatorData struct {
	pricesGen *generator.PricesGenerator
}

func NewGeneratorData(gen *generator.PricesGenerator) PricesData {
	return &generatorData{pricesGen: gen}
}

func (gp *generatorData) GetPrices(ctx context.Context) <-chan domain.CandleData {
	in := gp.pricesGen.Prices(ctx)
	out := make(chan domain.CandleData)

	go func() {
		defer close(out)
		for val := range in {
			out <- val
		}
	}()

	return out
}
