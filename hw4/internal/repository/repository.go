package repository

import "github.com/keruch/tfs-go-hw/hw4/internal/domain"

type Repo interface {
	GetUser(username string) (domain.UserData, error)
	GetAllUsers() []domain.UserData
	GetMessages() []domain.Message
	GetNumMessages(num int) ([]domain.Message, error)
	GetPrivateMessages(username string) ([]domain.Message, error)
}
