package utils

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/keruch/tfs-go-hw/trading_robot/internal/domain"
)

func ValidateDataIsPrice(data []byte) (domain.Price, bool) {
	var price domain.Price
	if err := json.Unmarshal(data, &price); err != nil {
		return domain.Price{}, false
	}

	if err := validator.New().Struct(price); err != nil {
		return domain.Price{}, false
	}

	return price, true
}

func ValidateDataIsTicker(data []byte) (domain.Ticker, bool) {
	var ticker domain.Ticker
	if err := json.Unmarshal(data, &ticker); err != nil {
		return domain.Ticker{}, false
	}

	if err := validator.New().Struct(ticker); err != nil {
		return domain.Ticker{}, false
	}

	return ticker, true
}
