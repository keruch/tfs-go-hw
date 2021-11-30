package tg

import (
	"context"
	"sync"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/keruch/tfs-go-hw/trading_robot/pkg/log"
)

const (
	startMsg = "/start"
	stopMsg  = "/stop"
)

type TelegramBot struct {
	bot *tgbot.BotAPI

	mu    sync.RWMutex     // mutex for protecting map
	users map[string]int64 // map with users

	logger *log.Logger
}

func NewTelegramBot(token string, logger *log.Logger) (*TelegramBot, error) {
	bot, err := tgbot.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &TelegramBot{
		bot:   bot,
		users: make(map[string]int64),
		logger: logger,
	}, nil
}

func (tg *TelegramBot) Serve(ctx context.Context) {
	u := tgbot.NewUpdate(0)
	u.Timeout = 60

	updates := tg.bot.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			tg.logger.Info("Telegram bot: serve done")
			return
		case update := <-updates:
			if update.Message == nil {
				continue
			}

			if update.Message.Text == startMsg {
				tg.addUser(update.Message.From.UserName, update.Message.Chat.ID)
			} else if update.Message.Text == stopMsg {
				tg.removeUser(update.Message.From.UserName)
			}
		}
	}

}

func (tg *TelegramBot) NotifyUsers(message string) {
	tg.mu.RLock()
	for user, ID := range tg.users {
		msg := tgbot.NewMessage(ID, message)

		_, err := tg.bot.Send(msg)
		if err != nil {
			tg.logger.Errorf("Send msg to %s user failed: %s", user, err)
		}
	}
	tg.mu.RUnlock()
}

func (tg *TelegramBot) addUser(username string, chatID int64) {
	tg.mu.Lock()
	tg.logger.Debugf("Added new user %s", username)
	tg.users[username] = chatID
	tg.mu.Unlock()
}

func (tg *TelegramBot) removeUser(username string) {
	tg.mu.Lock()
	tg.logger.Debugf("Removed user %s", username)
	delete(tg.users, username)
	tg.mu.Unlock()
}
