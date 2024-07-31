package services

import (
	"context"
	pb "github.com/megotours/proto/auth"
	"google.golang.org/grpc"
	"log"
	"time"
)

func Connect() {
	conn, err := grpc.NewClient("localhost:50051")
	if err != nil {
		log.Fatalf("Не удалось подключиться: %v", err)
	}
	defer conn.Close()

	client := pb.NewUsersClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.Get(ctx, &pb.GetByIdRequest{Id: "1"})
	if err != nil {
		log.Fatalf("Не удалось получить данные: %v", err)
	}

	log.Printf("Ответ: %s", response.GetUser())
}
