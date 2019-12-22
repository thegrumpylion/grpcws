package grpcws

import (
	"context"
	"net"
	"nhooyr.io/websocket"
	"time"
)

func NewDialer(ctx context.Context, opts *websocket.DialOptions) func(s string, dt time.Duration) (net.Conn, error) {
	return func(s string, dt time.Duration) (net.Conn, error) {
		c, _, err := websocket.Dial(ctx, s, opts)
		if err != nil {
			return nil, err
		}
		return websocket.NetConn(ctx, c, websocket.MessageBinary), nil
	}
}
