package main

import (
	"context"
	"fmt"
	userpb "go_project/service"
	U "go_project/utils"
	"sync"
	"time"

	"google.golang.org/grpc"
)

const (
	address_to_server   = "localhost:50051"
	address_to_listener = "localhost:50052"
)

// 구조체
type UserInfo struct {
	id   string
	name string
}

func InputData(u *UserInfo) {
	fmt.Printf("ID: ")
	fmt.Scanln(&u.id)

	fmt.Printf("Name: ")
	fmt.Scanln(&u.name)
}

func inInsertData(wg *sync.WaitGroup, ch chan int32, u *UserInfo) {
	conn, err := grpc.Dial(address_to_server, grpc.WithInsecure(), grpc.WithBlock())
	U.CheckErr(err)

	defer conn.Close()
	client := userpb.NewServiceInterfaceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	Result, err := client.InsertData(ctx, &userpb.InsertMsg{UserId: u.id, Name: u.name})
	U.CheckErr(err)

	ch <- Result.Index
	wg.Done()
}

func inListener(wg *sync.WaitGroup, ch chan int32) {
	n := <-ch

	conn, err := grpc.Dial(address_to_listener, grpc.WithInsecure(), grpc.WithBlock())
	U.CheckErr(err)

	defer conn.Close()
	client := userpb.NewServiceInterfaceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	Result, err := client.InsertIndex(ctx, &userpb.BroadcastIndex{Index: n})
	U.CheckErr(err)

	if Result.StatusCode == 200 {
		fmt.Printf("Response Code: %d, Success!\n", Result.StatusCode)
	}

	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	ch := make(chan int32)

	u := UserInfo{}
	InputData(&u)

	wg.Add(2)
	go inInsertData(&wg, ch, &u)
	go inListener(&wg, ch)

	wg.Wait()
}
