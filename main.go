package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	pb "qiwi/protobuf"
)

type Server struct{}

func main() {
	lis, err := net.Listen("tcp", ":"+Config.RPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterQiwiServer(s, &Server{})
	log.Println("Listen on " + Config.RPCPort + " port")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
