package main

import (
	"context"
	"flag"
	"io"
	"log"

	kis "github.com/uncommented/pfm/kis"
	upbit "github.com/uncommented/pfm/upbit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var serverAddr string
var kisAccountNumber string

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	flag.StringVar(&serverAddr, "server", "localhost:61000", "The address of grpc server")
	flag.StringVar(&kisAccountNumber, "kis", "", "The account number for KIS")
	flag.Parse()

	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Fail to connect %s: %v", serverAddr, err)
	}
	defer conn.Close()

	log.Println("----- Investments in KIS Account -----")
	kisStub := kis.NewKISAccountClient(conn)
	kisInvestmentStream, err := kisStub.ListInvestments(
		context.Background(),
		&kis.KISAccountRequest{
			AccountNumber: kisAccountNumber,
			MarketCode:    kis.MarketCode_NASD,
			Currency:      kis.Currency_USD,
		},
	)
	if err != nil {
		log.Println("Fail to call ListInvestments (KIS)")
	}

	i := 1
	for {
		kisInvestment, err := kisInvestmentStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListInvestments(_) = _, %v", kisStub, err)
		}
		log.Printf("Investment %d: %v\n", i, kisInvestment)
		i += 1
	}

	log.Println("----- Performance of KIS Account -----")
	performance, err := kisStub.GetPerformance(
		context.Background(),
		&kis.KISAccountRequest{
			AccountNumber: kisAccountNumber,
			MarketCode:    kis.MarketCode_NASD,
			Currency:      kis.Currency_USD,
		},
	)
	if err != nil {
		log.Println("Fail to call GerPerformance (KIS)")
	}
	log.Printf("Performance: %v\n", performance)

	log.Println("----- Investments of Upbit Account -----")
	upbitStub := upbit.NewUpbitAccountClient(conn)
	upbitInvestmentStream, err := upbitStub.ListInvestments(
		context.Background(),
		&upbit.UpbitAccountRequest{},
	)
	if err != nil {
		log.Println("Fail to call ListInvestments (Upbit)")
	}
	i = 1
	for {
		upbitInvestment, err := upbitInvestmentStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListInvestments(_) = _, %v", upbitStub, err)
		}
		log.Printf("Investment %d: %v\n", i, upbitInvestment)
		i += 1
	}
}
