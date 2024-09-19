package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	pb "filesharep2p/fileshare"
	reg "filesharep2p/register"
	sum "filesharep2p/sum"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port = flag.Int("port", 50052, "The server port")
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

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		// Verifica se é do tipo IP
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			// Pega apenas endereços IPv4
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}

	return ""
}

func calculateSum() {
	c := make(chan sum.Sums)
	nt := sum.ReadFiles(c)
	for i := 0; i < nt; i++ {
		x := <-c
		sums[x.Sum] = append(sums[x.Sum], x.Path)
	}
}

func initiateServer() {
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

func register() {
	ip := getLocalIP()
	connReg, err := grpc.NewClient(ip+":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Falha ao conectar ao servidor: %v", err)
	}
	defer connReg.Close()
	client := reg.NewRegisterServiceClient(connReg)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	registerResp, err := client.Register(ctx, &reg.RegisterRequest{Address: fmt.Sprintf("%s:%d", ip, *port)})
	if err != nil {
		log.Fatalf("Erro ao registrar o peer: %v", err)
	}
	fmt.Print(registerResp.Message)
}

func main() {
	flag.Parse()
	go calculateSum()
	go register()
	initiateServer()
}
