package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/uncommented/pfm/kis"
	pb "github.com/uncommented/pfm/portfolio"
	"github.com/uncommented/pfm/upbit"
	"google.golang.org/grpc"
)

type portfolioServer struct {
	pb.UnimplementedPortfolioServer
}

func (ps *portfolioServer) GetBalance(balanceRequest *pb.BalanceRequest, stream pb.Portfolio_GetBalanceServer) error {
	accountNumber := balanceRequest.AccountNumber
	vendor := balanceRequest.Vendor
	currency := balanceRequest.Currency

	if strings.EqualFold(vendor, "kis") {
		kisBalance := kis.RequestBalance(accountNumber, currency)
		balancesPerItem := kisBalance["output1"].([]interface{})

		for _, _balance := range balancesPerItem {
			balance := _balance.(map[string]interface{})
			itemType, _ := strconv.ParseInt(balance["prdt_type_cd"].(string), 10, 64)
			quantity, _ := strconv.ParseFloat(balance["ord_psbl_qty"].(string), 64)
			purchsingPrice, _ := strconv.ParseFloat(balance["pchs_avg_pric"].(string), 64)
			currentPrice, _ := strconv.ParseFloat(balance["now_pric2"].(string), 64)
			profitLossRate, _ := strconv.ParseFloat(balance["evlu_pfls_rt"].(string), 64)
			item := pb.BalanceItem{
				Type:            itemType,
				Name:            balance["ovrs_item_name"].(string),
				Ticker:          balance["ovrs_pdno"].(string),
				Quantity:        quantity,
				PurchasingPrice: purchsingPrice,
				CurrentPrice:    currentPrice,
				ProfitLossRate:  profitLossRate,
			}

			if err := stream.Send(&item); err != nil {
				return err
			}
		}
		return nil
	} else {
		balancesPerItem := upbit.RequestBalance()

		for _, balance := range balancesPerItem {
			quantity, _ := strconv.ParseFloat(balance["balance"].(string), 64)
			purchasingPrice, _ := strconv.ParseFloat(balance["avg_buy_price"].(string), 64)
			currency := balance["currency"].(string)

			if currency == "KRW" {
				continue
			}

			marketSnapshot := upbit.RequestMarketSnapshot(currency)
			currentPrice := marketSnapshot["trade_price"].(float64)
			profitLossRate := (currentPrice - purchasingPrice) / purchasingPrice

			marketInfo := upbit.RequestMarketInfo(currency)
			name := marketInfo["english_name"].(string)

			item := pb.BalanceItem{
				Type:            999,
				Name:            name,
				Ticker:          currency,
				Quantity:        quantity,
				PurchasingPrice: purchasingPrice,
				CurrentPrice:    currentPrice,
				ProfitLossRate:  profitLossRate,
			}

			if err := stream.Send(&item); err != nil {
				return err
			}
		}
		return nil
	}
}

func (ps *portfolioServer) GetPerformance(ctx context.Context, balanceRequest *pb.BalanceRequest) (*pb.Performance, error) {
	AccountNumber := balanceRequest.AccountNumber
	Vendor := balanceRequest.Vendor
	Currency := balanceRequest.Currency

	if strings.EqualFold(Vendor, "kis") {
		kisBalance := kis.RequestBalance(AccountNumber, Currency)
		overallPerformance := kisBalance["output2"].(map[string]interface{})
		totalPurchasingAmount, _ := strconv.ParseFloat(overallPerformance["frcr_pchs_amt1"].(string), 64)
		totalEvaluationAmount, _ := strconv.ParseFloat(overallPerformance["tot_evlu_pfls_amt"].(string), 64)
		totalProfitLoss, _ := strconv.ParseFloat(overallPerformance["ovrs_tot_pfls"].(string), 64)
		totalProfitLossRate, _ := strconv.ParseFloat(overallPerformance["tot_pftrt"].(string), 64)
		performance := pb.Performance{
			TotalPurchasingAmount: totalPurchasingAmount,
			TotalEvaluationAmount: totalEvaluationAmount,
			TotalProfitLoss:       totalProfitLoss,
			TotalProfitLossRate:   totalProfitLossRate,
		}
		return &performance, nil
	} else {
		balancesPerItem := upbit.RequestBalance()

		totalPurchasingAmount := 0.0
		totalEvaluationAmount := 0.0

		for _, balance := range balancesPerItem {
			quantity, _ := strconv.ParseFloat(balance["balance"].(string), 64)
			purchasingPrice, _ := strconv.ParseFloat(balance["avg_buy_price"].(string), 64)
			currency := balance["currency"].(string)
			if currency == "KRW" {
				continue
			}
			marketSnapshot := upbit.RequestMarketSnapshot(currency)
			currentPrice := marketSnapshot["trade_price"].(float64)

			totalPurchasingAmount += purchasingPrice * quantity
			totalEvaluationAmount += currentPrice * quantity
		}
		totalProfitLoss := totalEvaluationAmount - totalPurchasingAmount
		totalProfitLossRate := totalProfitLoss / totalPurchasingAmount
		performance := pb.Performance{
			TotalPurchasingAmount: totalPurchasingAmount,
			TotalEvaluationAmount: totalEvaluationAmount,
			TotalProfitLoss:       totalProfitLoss,
			TotalProfitLossRate:   totalProfitLossRate,
		}
		return &performance, nil
	}
}

var port int

func newServer() *portfolioServer {
	s := &portfolioServer{}
	return s
}

func main() {
	flag.IntVar(&port, "port", 61000, "The server port")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterPortfolioServer(grpcServer, newServer())
	_ = grpcServer.Serve(lis)
}
