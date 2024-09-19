package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "filesharep2p/fileshare"
	sum "filesharep2p/sum"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
	sums = make(map[int][]string)
)

type server struct {
	pb.UnimplementedFileshareServiceServer
}

func (s *server) HasFile(ctx context.Context, in *pb.MessageRequest) (*pb.MessageResponse, error) {
	log.Printf("Received: %v", in.GetHash())
	if sums[int(in.GetHash())] != nil {
		return &pb.MessageResponse{HasFile: true}, nil
	}
	return &pb.MessageResponse{HasFile: false}, nil
}

func main() {
	c := make(chan sum.Sums)
	nt := sum.ReadFiles(c)
	for i := 0; i < nt; i++ {
		x := <-c
		sums[x.Sum] = append(sums[x.Sum], x.Path)
	}

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
