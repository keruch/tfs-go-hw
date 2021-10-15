package main

import (
	"os"

	"github.com/keruch/tfs-go-hw/hw4/internal/domain/storage"
	"github.com/keruch/tfs-go-hw/hw4/internal/services"
	"github.com/keruch/tfs-go-hw/hw4/pkg/log"
)

func init() {
	// For debugging. In general, it is better to store the secret somewhere, but not in the source :)
	err := os.Setenv("ACCESS_SECRET", "toss_a_coin_to_your_witcher")
	if err != nil {
		return
	}
}

func main() {
	logger := log.NewLogger()
	inMemoryStorage := storage.NewInMemoryStorage()
	service := services.NewChatService(inMemoryStorage, inMemoryStorage, logger)
	service.Start()
}
