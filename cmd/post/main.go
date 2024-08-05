package main

import (
	adapter "antibomberman/mego-post/internal/adapters/grpc"
	"antibomberman/mego-post/internal/config"
	"antibomberman/mego-post/internal/database"
	"antibomberman/mego-post/internal/repositories"
	"antibomberman/mego-post/internal/services"
	"context"
	userPb "github.com/antibomberman/mego-protos/gen/go/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"time"
)

func main() {
	cfg := config.Load()
	db, err := database.ConnectToDB(cfg)
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	postRepository := repositories.NewPostRepository(db)
	if err != nil {
		log.Fatal(err)
	}
	/////////////////////////////User Service Client Connectivity Check
	userConn, err := grpc.NewClient("mego_user:20241", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to user service: %v", err)
	}
	defer func(userConn *grpc.ClientConn) {
		err := userConn.Close()
		if err != nil {
			log.Printf("Warning: Failed to close user service connection: %v", err)
		}
	}(userConn)

	time.Sleep(2 * time.Second)
	userClient := userPb.NewUserServiceClient(userConn)

	userDetails, err := userClient.GetById(ctx, &userPb.Id{Id: "1"})
	if err != nil {
		log.Printf("Warning: Failed to get user details: %v", err)
	} else {
		log.Printf("User details: %+v", userDetails)
	}
	////////////////////////////////////////////////////////////////////////////////////
	postService := services.NewPostService(postRepository, grpc.WithTransportCredentials(insecure.NewCredentials()))

	l, err := net.Listen("tcp", ":"+cfg.PostServiceServerPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gRPC := grpc.NewServer()

	adapter.Register(gRPC, cfg, postService)

	if err := gRPC.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
