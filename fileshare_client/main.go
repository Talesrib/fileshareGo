package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "filesharep2p/fileshare"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	dhash = "1023"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	hash = flag.String("hash", dhash, "Hash to search")
)

func main() {
	flag.Parse()

	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewFileshareServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.HasFile(ctx, &pb.MessageRequest{Hash: *hash})
	if err != nil {
		log.Fatalf("could not search: %v", err)
	}
	if r.HasFile {
		log.Printf("founded!")
	}
}
