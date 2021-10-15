package handlers

import "github.com/keruch/tfs-go-hw/hw4/internal/domain"

type Controller interface {
	SaveUserData(data domain.UserData) error
	SaveMessage(message domain.Message)
	SendPrivateMessage(username string, message domain.Message)
}
