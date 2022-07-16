package main

import (
	"context"
	"fmt"
	userpb "go_project/service"
	U "go_project/utils"
	"log"
	"net"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
)

const (
	portnumber = "50051"
	mongo_URI  = "mongodb://localhost:27017"
)

var db *mongo.Client
var client *mongo.Collection
var index int32 = 0

type server struct {
	userpb.ServiceInterfaceServer
}

type UserInfo struct {
	UserId string `bson:"UserID"`
	Name   string `bson:"Name"`
}

// Server - Client gRPC Call
func (s *server) InsertData(ctx context.Context, req *userpb.InsertMsg) (*userpb.InsertResponse, error) {
	var S_Code int32 = 500
	UserId := req.UserId
	UserName := req.Name
	index++

	//findIndexinDB()

	Data := bson.D{
		{"Index", index},
		{"UserID", UserId},
		{"Name", UserName},
	}
	result, err := client.InsertOne(context.TODO(), Data)
	U.CheckErr(err)

	fmt.Println(result.InsertedID)
	S_Code -= 300
	return &userpb.InsertResponse{Index: index, StatusCode: S_Code}, nil
}

// Server - Listener gRPC Call
func (s *server) GetData(ctx context.Context, req *userpb.GetMsg) (*userpb.GetResponse, error) {
	IndexNum := req.Index

	data := UserInfo{}

	// 필터 조건 선언하기 : Index 필드의 값이 IndexNum인 document만 Read
	filter := bson.M{"Index": IndexNum}
	if err := client.FindOne(context.TODO(), filter).Decode(&data); err != nil {
		panic(err)
	}

	response := &userpb.GetResponse{
		UserId: data.UserId,
		Name:   data.Name,
	}

	return response, nil
}

func findIndexinDB() {
	var index primitive.ObjectID
	var result bson.M
	opts := options.FindOne().SetSort(bson.D{{"Index", 1}})

	err := client.FindOne(
		context.TODO(),
		bson.D{{"Index", index}},
		opts,
	).Decode(&result)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in
		// the collection.
		if err == mongo.ErrNoDocuments {
			return
		}
		log.Fatal(err)
	}

	fmt.Printf("found document %s", result)
}

// ServerEngine (Server - Listener, Client)
func ServerEngine() {
	// Server Connect
	lis, err := net.Listen("tcp", ":"+portnumber)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterServiceInterfaceServer(grpcServer, &server{})

	// MongoDB Connect
	fmt.Println("Connecting to MongoDB...")
	db, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongo_URI))
	U.CheckErr(err)

	if err := db.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Connect to MongoDB!!")
	client = db.Database("mydb").Collection("users")

	log.Printf("start gRPC server on %s port", portnumber)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to server: %s", err)
	}
}

func main() {
	ServerEngine()
}
