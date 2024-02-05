package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedGetInfoServer
}

func main() {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to mongoDB server: %v", err)
	}
	client := protobuf.NewMongoDBClient(conn)
	ctx := context.Background()
	db, err := client.GetDB(ctx, &protobuf.Empty{})
	if err != nil {
		log.Fatalf("error getting DB from server: %v", err)
	}

}
