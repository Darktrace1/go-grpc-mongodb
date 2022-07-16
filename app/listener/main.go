package main

import (
	"context"
	"fmt"
	userpb "go_project/service"
	U "go_project/utils"
	"log"
	"net"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
)

var db *mongo.Client
var client *mongo.Collection

const (
	portnumber        = "50052"
	address_to_server = "localhost:50051"
)

type server struct {
	userpb.ServiceInterfaceServer
}

// Server - Listener gRPC Call
func (s *server) InsertIndex(ctx context.Context, req *userpb.BroadcastIndex) (*userpb.BroadcastResponse, error) {
	index := req.Index

	SendIndex(index)

	return &userpb.BroadcastResponse{StatusCode: 200}, nil

}

// ServerEngine (Listener - Client)
func ServerEngine() {
	// server ON
	lis, err := net.Listen("tcp", ":"+portnumber)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterServiceInterfaceServer(grpcServer, &server{})

	// mongodb
	fmt.Println("Connecting to MongoDB...")
	db, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	if err := db.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Connect to MongoDB!!")

	client = db.Database("mydb").Collection("users")

	// server print
	log.Printf("start gRPC server on %s port", portnumber)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to server: %s", err)
	}
}

// Send Index to Server
func SendIndex(index int32) {
	conn, err := grpc.Dial(address_to_server, grpc.WithInsecure(), grpc.WithBlock())
	U.CheckErr(err)

	defer conn.Close()
	client := userpb.NewServiceInterfaceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.GetData(ctx, &userpb.GetMsg{Index: index})
	U.CheckErr(err)

	fmt.Printf("User Id: %s, ", r.UserId)
	fmt.Printf("User Name: %s\n", r.Name)
}

func main() {
	ServerEngine()
}
