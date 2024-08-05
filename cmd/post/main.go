package main

import (
	adapter "antibomberman/mego-post/internal/adapters/grpc"
	"antibomberman/mego-post/internal/clients"
	"antibomberman/mego-post/internal/config"
	"antibomberman/mego-post/internal/database"
	"antibomberman/mego-post/internal/repositories"
	"antibomberman/mego-post/internal/services"
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
	storageClient, err := clients.NewStorageClient(cfg.UserServiceAddress)
	postService := services.NewPostService(
		postRepository,
		postContentRepository,
		postContentFileRepository,
		userClient,
		storageClient,
	)

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
