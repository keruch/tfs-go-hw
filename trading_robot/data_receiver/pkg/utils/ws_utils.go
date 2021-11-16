package utils

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

var ErrBadHandshake = errors.New("can not handshake with host")

func CreateRetryableDial(u url.URL, maxRetries int) (*websocket.Conn, *http.Response, error) {
	var (
		conn         *websocket.Conn
		resp         *http.Response
		retriesCount int
	)
	for {
		var err error
		conn, resp, err = websocket.DefaultDialer.Dial(u.String(), nil)
		if err == nil {
			break
		}
		if errors.Is(err, websocket.ErrBadHandshake) {
			retriesCount++
			if retriesCount > maxRetries {
				return nil, nil, ErrBadHandshake
			}
			continue
		}
		return nil, nil, err
	}

	return conn, resp, nil
}
