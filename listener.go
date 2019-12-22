// +build !js

package grpcws

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"nhooyr.io/websocket"
)

func NewListener(ctx context.Context, lis net.Listener, srv *http.Server, mux *http.ServeMux, path string) net.Listener {
	wsl := &wsListener{
		stop: make(chan struct{}),
		errc: make(chan error, 1),
		conn: make(chan net.Conn),
	}

	if path == "" {
		path = "/ws"
	}

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		done := make(chan struct{})
		c := &wsConn{
			Conn: websocket.NetConn(ctx, conn, websocket.MessageBinary),
			done: done,
		}
		wsl.conn <- c
		<-done
	})

	srv.Handler = mux

	go func() {
		defer close(wsl.errc)
		wsl.errc <- srv.Serve(lis)
	}()

	return wsl
}

type wsAddr struct {
}

func (wsAddr) Network() string {
	return "websocket"
}

func (wsAddr) String() string {
	return "websocket/unknown-addr"
}

type wsConn struct {
	net.Conn
	done chan struct{}
	once sync.Once
}

func (c *wsConn) Close() error {
	c.once.Do(func() {
		close(c.done)
	})
	return c.Conn.Close()
}

type wsListener struct {
	stop chan struct{}
	errc chan error
	conn chan net.Conn
}

func (l *wsListener) Accept() (net.Conn, error) {
	select {
	case <-l.stop:
		return nil, fmt.Errorf("server stopped")
	case err := <-l.errc:
		l.Close()
		return nil, err
	case c := <-l.conn:
		return c, nil
	}
}

func (l *wsListener) Close() error {
	select {
	case <-l.stop:
	default:
		close(l.stop)
	}
	return nil
}

func (l *wsListener) Addr() net.Addr {
	return wsAddr{}
}
