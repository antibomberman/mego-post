package main

import (
	adapter "antibomberman/mego-post/internal/adapters/grpc"
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
	defer db.Close()

	postRepository := repositories.NewPostRepository(db)
	if err != nil {
		log.Fatal(err)
	}
	postService := services.NewPostService(postRepository)

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
