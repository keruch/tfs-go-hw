package utils

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

var ErrBadHandshake = errors.New("can not handshake with host")

// RetryableWSConn - should be one instance for each dial
type RetryableWSConn struct {
	Url        url.URL
	MaxRetries int
	RequestHeader http.Header

	conn       *websocket.Conn
}

func (c *RetryableWSConn) RetryableDial() (*http.Response, error) {
	var (
		conn         *websocket.Conn
		resp         *http.Response
		retriesCount int
	)
	for {
		var err error
		conn, resp, err = websocket.DefaultDialer.Dial(c.Url.String(), c.RequestHeader)
		if err == nil {
			break
		}
		if errors.Is(err, websocket.ErrBadHandshake) {
			retriesCount++
			if retriesCount > c.MaxRetries {
				return nil, ErrBadHandshake
			}
			continue
		}
		return nil, err
	}

	c.conn = conn
	return resp, nil
}

func (c *RetryableWSConn) ReadMessage() (messageType int, p []byte, reconnected bool, err error) {
	messageType, p, err = c.conn.ReadMessage()
	if err != nil {
		if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			_, err = c.RetryableDial()
			if err != nil {
				return
			}
			reconnected = true
			return
		}
		return
	}
	return
}

func (c *RetryableWSConn) WriteJSON(v interface{}) (reconnected bool, err error) {
	err = c.conn.WriteJSON(v)
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			_, err = c.RetryableDial()
			if err != nil {
				return
			}
			reconnected = true
			return
		}
		return
	}
	return
}

func (c *RetryableWSConn) Close() error {
	return c.conn.Close()
}
