package main

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/thegrumpylion/grpcws"
	"github.com/thegrumpylion/grpcws/example/service"
	"google.golang.org/grpc"
)

func main() {

	ctx := context.Background()

	conn, err := grpc.Dial("ws://localhost:8085/ws",
		grpc.WithDialer(grpcws.NewDialer(ctx, nil)),
		grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	cli := service.NewTrackerClient(conn)

	// normal call

	resp, err := cli.Login(ctx, &service.LoginReq{
		Username: "tester",
		Password: "test",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("resp", resp)

	// server stream

	stream, err := cli.Events(ctx, &service.EventsReq{})
	if err != nil {
		panic(err)
	}

	for {
		ev, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Println("ev:", ev.Event)
	}

	// client stream

	stream2, err := cli.Track(ctx)
	if err != nil {
		panic(err)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	n := (r.Int() % 5) + 5

	fmt.Println("will send:", n)

	for i := 0; i < n; i++ {
		if err := stream2.Send(&service.TrackReq{
			Lng: r.ExpFloat64(),
			Lat: r.ExpFloat64(),
		}); err != nil {
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
		}
	}
	reply, err := stream2.CloseAndRecv()
	fmt.Println("track count", reply.Count)

	select {}
}
