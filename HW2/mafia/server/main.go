package main

import (
	"log"
	"mafia/internal/game"
	"net"

	"google.golang.org/grpc"
)

func main() {
	// create listiner
	lis, err := net.Listen("tcp", ":50005")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create grpc server
	s := grpc.NewServer()
	game.RegisterConnectToGameServer(s, server{})

	log.Println("start server")
	// and start...
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type server struct{}

func (s server) ConnectToGame(req *game.JoinGameRoomRequest, stream game.ConnectToGame_DoServer) error {

	return nil
}
