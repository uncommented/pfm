package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/joho/godotenv"
	kis "github.com/uncommented/pfm/portfolio/kis"
	upbit "github.com/uncommented/pfm/portfolio/upbit"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env")
		os.Exit(1)
	}
	log.SetFlags(log.LstdFlags | log.Llongfile)

	lis, err := net.Listen("tcp", ":61000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		os.Exit(1)
	}
	s := grpc.NewServer()

	kis.RegisterKISAccountServer(s, &kis.KISAccountService{})
	upbit.RegisterUpbitAccountServer(s, &upbit.UpbitAccountService{})
	log.Printf("Service listening at %v", lis.Addr())

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		os.Exit(1)
	}
}
