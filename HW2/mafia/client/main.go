package main

import (
	"context"
	"fmt"
	"io"

	"log"
	"mafia/internal/game"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func main() {
	conn, err := grpc.Dial(":50005", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}

	// create stream
	client := game.NewConnectToGameClient(conn)

	request := &game.JoinGameRoomRequest{
		Name: "nadya",
	}
	stream, err := client.Do(context.Background(), request)

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}

	resp, err := stream.Recv()
	roomNumb := resp.RoomNumb

	if err == io.EOF {
		//done <- true //means stream is finished
		return
	}

	fmt.Println("Game Start, your room numb is", roomNumb)

	//fmt.Println(response.Message)
}
