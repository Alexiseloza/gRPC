package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"

	pb "grpcserver/grpcServer/proto-grpc"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedGetInfoServer
}

func save_comments(comment string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoclient, err := mongo.Connect(ctx, options.Client().ApplyURI("momgo uri"))
	if err != nil {
		log.Fatal(err)
	}

	databases, err := mongoclient.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)

	Database := mongoclient.Database("testGRPC")
	Collection := Database.Collection("comments")

	var bdoc interface{}

	errb := bson.UnmarshalExtJSON([]byte(comment), true, &bdoc)
	fmt.Println(errb)

	insertResult, err := Collection.InsertOne(ctx, bdoc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insertResult)
}

func (s *server) ReturnInfo(ctx context.Context, in *pb.RequestId) (*pb.ReplyInfo, error) {
	save_comments(in.GetId())
	fmt.Printf(">> We have obtaind the data client: %n", in.GetId())
	return &pb.ReplyInfo{Info: ">> Hi! thanks for you comment: " + in.GetId()}, nil
}

func main() {
	hear, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("We can't do start the server: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGetInfoServer(s, &server{})
	if err := s.Serve(hear); err != nil {
		log.Fatal("We can't do start the server: %v", err)
	}
	fmt.Println(port + " is listening...")

}
