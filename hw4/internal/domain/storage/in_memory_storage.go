package storage

import (
	"errors"
	"sync"

	"github.com/keruch/tfs-go-hw/hw4/internal/domain"
)

var (
	ErrorUserAlreadyExist = errors.New("user with this name already exists")
	ErrorUserNotFound     = errors.New("user not found")
	ErrorMessageboxEmpty  = errors.New("messagebox is empty")
)

type (
	usersTable struct {
		mutex sync.RWMutex
		table map[string]domain.UserData
	}
	msgArray struct {
		mutex sync.RWMutex
		array []domain.Message
	}
	privateMsgTable struct {
		mutex sync.RWMutex
		table map[string][]domain.Message
	}
)

func newUsersTable() usersTable {
	return usersTable{
		table: make(map[string]domain.UserData),
	}
}

func newMsgArray() msgArray {
	return msgArray{
		array: make([]domain.Message, 0),
	}
}

func newPrivateMsgTable() privateMsgTable {
	return privateMsgTable{
		table: make(map[string][]domain.Message),
	}
}

type InMemoryStorage struct {
	u usersTable      // users
	m msgArray        // messages
	p privateMsgTable // private messages
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		u: newUsersTable(),
		m: newMsgArray(),
		p: newPrivateMsgTable(),
	}
}

func (ms *InMemoryStorage) GetUser(username string) (domain.UserData, error) {
	ms.u.mutex.RLock()
	user, ok := ms.u.table[username]
	ms.u.mutex.RUnlock()

	if !ok {
		return domain.UserData{}, ErrorUserNotFound
	}
	return user, nil
}

func (ms *InMemoryStorage) GetAllUsers() []domain.UserData {
	users := make([]domain.UserData, 0)

	ms.u.mutex.RLock()
	for _, val := range ms.u.table {
		users = append(users, val)
	}
	ms.u.mutex.RUnlock()

	return users
}

func (ms *InMemoryStorage) GetMessages() []domain.Message {
	tempArray := make([]domain.Message, 0)

	ms.m.mutex.RLock()
	tempArray = append(tempArray, ms.m.array...)
	ms.m.mutex.RUnlock()

	return tempArray
}

func (ms *InMemoryStorage) GetNumMessages(num int) ([]domain.Message, error) {
	count := len(ms.m.array) - num
	if count < 0 {
		count = 0
	}

	tempArray := make([]domain.Message, 0)

	ms.m.mutex.RLock()
	tempArray = append(tempArray, ms.m.array[count:]...)
	ms.m.mutex.RUnlock()

	return tempArray, nil
}

func (ms *InMemoryStorage) GetPrivateMessages(username string) ([]domain.Message, error) {
	ms.p.mutex.RLock()
	messages, ok := ms.p.table[username]
	ms.p.mutex.RUnlock()

	if !ok {
		return nil, ErrorMessageboxEmpty
	}
	return messages, nil
}

func (ms *InMemoryStorage) SaveUserData(data domain.UserData) error {
	ms.u.mutex.RLock()
	_, ok := ms.u.table[data.Username]
	ms.u.mutex.RUnlock()

	if ok {
		return ErrorUserAlreadyExist
	}

	ms.u.mutex.Lock()
	ms.u.table[data.Username] = data
	ms.u.mutex.Unlock()

	return nil
}

func (ms *InMemoryStorage) SaveMessage(message domain.Message) {
	ms.m.mutex.Lock()
	ms.m.array = append(ms.m.array, message)
	ms.m.mutex.Unlock()
}

func (ms *InMemoryStorage) SendPrivateMessage(username string, message domain.Message) {
	ms.p.mutex.RLock()
	_, ok := ms.u.table[username]
	ms.p.mutex.RUnlock()

	ms.p.mutex.Lock()
	if !ok {
		ms.p.table[username] = make([]domain.Message, 0)
	}

	ms.p.table[username] = append(ms.p.table[username], message)
	ms.p.mutex.Unlock()
}
