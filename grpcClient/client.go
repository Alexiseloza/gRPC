package main

import (
	"context"
	"encoding/json"
	"fmt"
	pb "grpcclient/grpcClient/proto-grpc"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func connect_server(writer http.ResponseWriter, req *http.Request) {
	//CORS
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	writer.Header().Set("Content-Type", "application/json")
	if req.Method == "GET" {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(`{"messagge":"Server is running"}`))
		return
	}

	datos, _ := io.ReadAll(req.Body)

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		json.NewEncoder(writer).Encode("Error connection to grpc Server")
		log.Fatal("Is not possible conncet ion with the server ", err)
	}
	defer conn.Close()

	cl := pb.NewGetInfoClient(conn)

	id := string(datos)
	if len(os.Args) > 1 {
		id = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	ret, err := cl.ReturnInfo(ctx, &pb.RequestId{Id: id})
	if err != nil {
		json.NewEncoder(writer).Encode("There was an error while connecting with the server")
		log.Fatal("Inposible catch the info", err)
	}
	log.Printf("Server response: %s\n", ret.GetMessage())
	json.NewEncoder(writer).Encode("Info saved")

}

func main() {
	router := mux.NewRouter().StrictSlach(true)
	router.HandleFunc("/", connect_server)
	fmt.Println("Client running on port : 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
