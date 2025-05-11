package main

import (
	"log"
	"net"

	"github.com/Vladimir220/subpub/subpub_service"
	"google.golang.org/grpc"
)

func listening(myServ *Server, url string, infoLog, errLog *log.Logger) {
	infoLog.Println("Registering a Subpub server by URL:", url)

	lis, err := net.Listen("tcp", url)
	if err != nil {
		errLog.Panicln(err)
	}

	var opts []grpc.ServerOption
	server := grpc.NewServer(opts...)

	subpub_service.RegisterPubSubServer(server, myServ)
	err = server.Serve(lis)
	if err != nil {
		errLog.Panicln(err)
	}

	infoLog.Println("Server stopped successfully")
}
