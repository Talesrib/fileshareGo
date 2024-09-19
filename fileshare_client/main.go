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
	port = flag.Int("port", 60001, "the address to connect to")
	sums = make([]sum.Sums, 0)
)

func calculateSum(done chan bool) {
	c := make(chan sum.Sums)
	nt := sum.ReadFiles(c)
	for i := 0; i < nt; i++ {
		x := <-c
		sums = append(sums, sum.Sums{Path: x.Path, Sum: x.Sum})
	}
	done <- true
}

func connect(address string, done chan bool) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		done <- true
		return
	}
	defer conn.Close()
	c := pb.NewFileshareServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for _, s := range sums {
		r, err := c.HasFile(ctx, &pb.MessageRequest{Hash: int64(s.Sum)})
		if err != nil {
			done <- true
			return
		}
		if r.HasFile {
			log.Printf("founded %s in %s", s.Path, address)
		}
	}
	done <- true
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

func getListOfPeers() *reg.ListOfPeersResponse {
	ip := getLocalIP()
	conn, err := grpc.NewClient(ip+":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Falha ao conectar ao servidor: %v", err)
	}
	defer conn.Close()
	client := reg.NewRegisterServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	localAddress := fmt.Sprintf("%s:%d", ip, *port)
	listOfPeers, err := client.GetPeers(ctx, &reg.ListOfPeersRequest{Address: localAddress})
	return listOfPeers
}

func main() {
	flag.Parse()
	done := make(chan bool)
	go calculateSum(done)
	listOfPeers := getListOfPeers()
	<-done
	nt := 0
	for _, peer := range listOfPeers.Peers {
		go connect(peer, done)
		nt += 1
	}
	for i := 0; i < nt; i++ {
		<-done
	}
}
