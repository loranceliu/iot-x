package main

import (
	"fmt"
	"google.golang.org/grpc"
	"iot-x/protobuf"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Printf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	protobuf.RegisterInstanceServer(s, &protobuf.Service{})

	if err := s.Serve(lis); err != nil {
		fmt.Printf("Failed to serve: %v", err)
	}
}
