package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/Vladimir220/subpub/subpub_service"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	err := godotenv.Load("./client/.env")
	if err != nil {
		panic(err)
	}
	url := os.Getenv("SERVER_URL")

	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panicln(err)
	}
	defer conn.Close()

	client := subpub_service.NewPubSubClient(conn)

	ctx, c := context.WithCancel(context.Background())
	defer c()

	go func() {
		time.Sleep(time.Second * 5)
		client.Publish(ctx, &subpub_service.PublishRequest{Key: "Огонь", Data: "Это же феерверк!"})
		time.Sleep(time.Second * 1)
		client.Publish(ctx, &subpub_service.PublishRequest{Key: "Огонь", Data: "Ого!"})
		time.Sleep(time.Second * 1)
		client.Publish(ctx, &subpub_service.PublishRequest{Key: "Огонь", Data: "Хыхы!"})
	}()

	stream, err := client.Subscribe(ctx, &subpub_service.SubscribeRequest{Key: "Огонь"})
	if err != nil {
		fmt.Println(err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			stream.CloseSend()
			fmt.Println(err)
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println("Получено:", msg.Data)
	}
}
