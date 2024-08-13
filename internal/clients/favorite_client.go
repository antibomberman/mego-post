package clients

import (
	pb "github.com/antibomberman/mego-protos/gen/go/favorite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type FavoriteClient struct {
	pb.FavoriteServiceClient
}

func NewFavoriteClient(address string) (*FavoriteClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &FavoriteClient{pb.NewFavoriteServiceClient(conn)}, nil
}
