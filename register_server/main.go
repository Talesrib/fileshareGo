package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "filesharep2p/register"

	"google.golang.org/grpc"
)

type bootstrapServer struct {
	pb.UnimplementedRegisterServiceServer
	mu    sync.Mutex
	peers []string
}

func (s *bootstrapServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	peerAddress := req.Address
	if peerAddress == "" {
		return &pb.RegisterResponse{
			Success: false,
			Message: "Invalid address",
		}, nil
	}

	// Registra o peer
	s.peers = append(s.peers, peerAddress)

	fmt.Printf("Peer %s registrado com sucesso\n", peerAddress)

	return &pb.RegisterResponse{
		Success: true,
		Message: "Peer registrado com sucesso",
	}, nil
}

func (s *bootstrapServer) GetPeers(ctx context.Context, req *pb.ListOfPeersRequest) (*pb.ListOfPeersResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	peers := []string{}
	for _, addr := range s.peers {
		if req.Address != addr {
			peers = append(peers, addr)
		}
	}

	return &pb.ListOfPeersResponse{
		Peers: peers,
	}, nil
}

func main() {
	server := &bootstrapServer{
		peers: make([]string, 0),
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Falha ao iniciar o listener: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRegisterServiceServer(grpcServer, server)

	fmt.Println("Servidor de Register rodando na porta 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
