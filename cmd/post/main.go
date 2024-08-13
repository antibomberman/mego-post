package main

import (
	adapter "antibomberman/mego-post/internal/adapters/grpc"
	"antibomberman/mego-post/internal/clients"
	"antibomberman/mego-post/internal/config"
	"antibomberman/mego-post/internal/database"
	"antibomberman/mego-post/internal/repositories"
	"antibomberman/mego-post/internal/services"
	"context"
	"fmt"
	pb "github.com/antibomberman/mego-protos/gen/go/storage"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	cfg := config.Load()
	db, err := database.ConnectToDB(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	postRepository := repositories.NewPostRepository(db)
	postContentRepository := repositories.NewPostContentRepository(db)
	postContentFileRepository := repositories.NewPostContentFileRepository(db)
	userClient, err := clients.NewUserClient(cfg.UserServiceAddress)
	storageClient, err := clients.NewStorageClient(cfg.StorageServiceAddress)
	favoriteClient, err := clients.NewFavoriteClient(cfg.StorageServiceAddress)

	postService := services.NewPostService(
		postRepository,
		postContentRepository,
		postContentFileRepository,
		userClient,
		storageClient,
		favoriteClient,
	)
	rsp, err := storageClient.GetObjectUrl(context.Background(), &pb.GetObjectUrlRequest{
		FileName: "test.png",
	})
	if err != nil {
		fmt.Printf("Error getting object URL: %v", err)
	}
	log.Println("Object URL:", rsp.Url)

	log.Printf("Starting gRPC server on port %s", cfg.PostServiceServerPort)
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
