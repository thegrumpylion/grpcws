package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/thegrumpylion/grpcws/example/service"
)

type svc struct {
	sync.Mutex
	users   map[string]string
	session map[string]time.Time
	events  []string
}

func getToken(user string) string {
	return user + ".AccessKey"
}

func (s *svc) Login(_ context.Context, req *service.LoginReq) (*service.Token, error) {
	s.Lock()
	defer s.Unlock()

	fmt.Println("login:", req)

	if pass, ok := s.users[req.Username]; ok {
		if req.Password == pass {
			tok := getToken(req.Username)
			s.session[tok] = time.Now().Add(time.Minute * 5)
			return &service.Token{Value: tok}, nil
		}
	}
	return nil, errors.New("unknown username/password")
}

func (s *svc) Logout(_ context.Context, tok *service.Token) (*service.LogoutRsp, error) {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.session[tok.Value]; ok {
		delete(s.session, tok.Value)
		return &service.LogoutRsp{Msg: "see ya"}, nil
	}
	return nil, errors.New("unknown token")
}

func (s *svc) Track(stream service.Tracker_TrackServer) error {
	var count uint32
	for {
		pos, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&service.TrackRsp{
				Count: count,
			})
		}
		if err != nil {
			return err
		}
		count++
		fmt.Println("pos:", pos)
	}
}

func (s *svc) Events(req *service.EventsReq, stream service.Tracker_EventsServer) error {
	for _, ev := range s.events {
		if err := stream.Send(&service.EventsRsp{Event: ev}); err != nil {
			return err
		}
	}
	return nil
}

func (s *svc) Chat(_ service.Tracker_ChatServer) error {
	panic("not implemented")
}

func NewService(users map[string]string) *svc {
	return &svc{
		users:   users,
		session: map[string]time.Time{},
		events:  mockEvents(),
	}
}

func mockEvents() []string {
	return []string{
		"info: adsfkl",
		"info: lkfdglsaj",
		"warning: ioaewjfo",
		"info: qweporid",
		"info: 903u4giog",
		"warning: fdfefeff",
		"info: mvkremie",
		"warning: ioaewjfo",
	}
}
