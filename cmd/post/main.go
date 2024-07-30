package post

import (
	adapter "antibomberman/mego-post/internal/adapters/grpc"
	"antibomberman/mego-post/internal/config"
	"antibomberman/mego-post/internal/repositories"
	"antibomberman/mego-post/internal/services"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	cfg := config.Load()
	db, err := sqlx.Open("postgres", cfg.DatabaseURL)
	err = db.Ping()
	defer db.Close()

	postRepository := repositories.NewPostRepository(db)
	if err != nil {
		log.Fatal(err)
	}
	srv := services.NewPostService(postRepository)

	l, err := net.Listen("tcp", ":"+cfg.ServerPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gRPC := grpc.NewServer()

	adapter.Register(gRPC, cfg, srv)

	if err := gRPC.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
