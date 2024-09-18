package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "filesharep2p/fileshare"
	sum "filesharep2p/sum"

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
	s := make(chan sum.Sums)
	nt := sum.ReadFiles(s)
	for i := 0; i < nt; i++ {
		x := <-s
		r, err := c.HasFile(ctx, &pb.MessageRequest{Hash: int64(x.Sum)})
		if err != nil {
			log.Fatalf("could not search: %v", err)
		}
		if r.HasFile {
			log.Printf("founded %s", x.Path)
		}
	}
}
