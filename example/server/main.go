package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/thegrumpylion/grpcws"
	"github.com/thegrumpylion/grpcws/example/service"
	"google.golang.org/grpc"
)

func main() {
	l, err := net.Listen("tcp", ":8085")
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(".")))

	lst := grpcws.NewListener(context.Background(), l, &http.Server{}, mux, "/ws")

	s := grpc.NewServer()

	srv := NewService(map[string]string{
		"user":   "password",
		"tester": "test",
	})

	service.RegisterTrackerServer(s, srv)

	fmt.Println("http://localhost:8085")

	err = s.Serve(lst)
	if err != nil {
		panic(err)
	}

}
