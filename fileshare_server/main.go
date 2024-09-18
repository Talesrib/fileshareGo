package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "filesharep2p/fileshare"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedFileshareServiceServer
}

func (s *server) HasFile(ctx context.Context, in *pb.MessageRequest) (*pb.MessageResponse, error) {
	log.Printf("Received: %v", in.GetHash())
	return &pb.MessageResponse{HasFile: true}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterFileshareServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
