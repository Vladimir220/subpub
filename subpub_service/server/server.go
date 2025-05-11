package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Vladimir220/subpub/subpub_lib"
	"github.com/Vladimir220/subpub/subpub_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	subpub_service.UnimplementedPubSubServer
	subpub     subpub_lib.SubPub
	stopSignal chan struct{}
	errLog     *log.Logger
}

func CreateServer(errLog *log.Logger) *Server {
	return &Server{subpub: subpub_lib.NewSubPub(), errLog: errLog, stopSignal: make(chan struct{})}
}

func (s *Server) Close() {
	close(s.stopSignal)
}

func (s *Server) Publish(ctx context.Context, req *subpub_service.PublishRequest) (*emptypb.Empty, error) {
	err := s.subpub.Publish(req.Key, req.Data)
	if err != nil {
		s.errLog.Println(err)
		return &emptypb.Empty{}, status.Error(codes.Internal, fmt.Sprint("Server error:", err))
	}
	return &emptypb.Empty{}, status.Error(codes.OK, "Success")
}

func (s *Server) Subscribe(req *subpub_service.SubscribeRequest, eventStream grpc.ServerStreamingServer[subpub_service.Event]) error {
	messageHandler := func(msg any) {
		res := &subpub_service.Event{
			Data: fmt.Sprint(msg),
		}
		eventStream.Send(res)
	}

	subscription, err := s.subpub.Subscribe(req.Key, messageHandler)
	if err != nil {
		s.errLog.Println(err)
		return status.Error(codes.Internal, fmt.Sprint("Server error:", err))
	}

	for err := s.subpub.Close(eventStream.Context()); err != nil; err = s.subpub.Close(eventStream.Context()) {
		select {
		case <-s.stopSignal:
			return status.Error(codes.Unavailable, "Stop signal. Server off.")
		default:
		}
	}

	subscription.Unsubscribe()

	return status.Error(codes.OK, "Subscription End.")
}
