package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	kis "github.com/uncommented/pfm/kis"
	upbit "github.com/uncommented/pfm/upbit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var conn *grpc.ClientConn
var upbitStub upbit.UpbitAccountClient
var kisStub kis.KISAccountClient

func getInvestments(c *gin.Context) {
	assetType := c.Query("assetType")

	var investments []Investment
	if strings.EqualFold(assetType, "security") {
		if kisStub == nil {
			kisStub = kis.NewKISAccountClient(conn)
		}
		kisInvestmentStream, err := kisStub.ListInvestments(
			context.Background(),
			&kis.KISAccountRequest{
				AccountNumber: kisAccountNumber,
				MarketCode:    kis.MarketCode_NASD,
				Currency:      kis.Currency_USD,
			},
		)
		if err != nil {
			InternalServerError(c,
				fmt.Sprintf("Failed to get investments (assetType: %s)", assetType))
		}
		for {
			kisInvestment, err := kisInvestmentStream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				InternalServerError(c,
					fmt.Sprintf("Failed to get investments (assetType: %s)", assetType))
			}
			_ = append(investments, Investment{
				ID:                     kisInvestment.SecurityCode,
				Name:                   kisInvestment.SecurityFullname,
				Quantity:               float64(kisInvestment.Quantity),
				AveragePurchasingPrice: kisInvestment.AveragePurchasingPrice,
				CurrentPrice:           kisInvestment.CurrentPrice,
				EvaluationAmount:       kisInvestment.EvaluationAmount,
				ProfitLoss:             kisInvestment.ProfitLoss,
				ProfitLossRate:         kisInvestment.ProfitLossRate,
			})
		}
	} else if strings.EqualFold(assetType, "crypto") {
		if upbitStub == nil {
			upbitStub = upbit.NewUpbitAccountClient(conn)
		}
		upbitInvestmentStream, err := upbitStub.ListInvestments(
			context.Background(),
			&upbit.UpbitAccountRequest{},
		)
		if err != nil {
			InternalServerError(c,
				fmt.Sprintf("Failed to get investments (assetType: %s)", assetType))
		}
		for {
			upbitInvestment, err := upbitInvestmentStream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				InternalServerError(c,
					fmt.Sprintf("Failed to get investments (assetType: %s)", assetType))
			}
			_ = append(investments, Investment{
				ID:                     upbitInvestment.Currency,
				Name:                   upbitInvestment.CurrencyFullname,
				Quantity:               upbitInvestment.Quantity,
				AveragePurchasingPrice: upbitInvestment.AveragePurchasingPrice,
				CurrentPrice:           upbitInvestment.CurrentPrice,
				EvaluationAmount:       upbitInvestment.EvaluationAmount,
				ProfitLoss:             upbitInvestment.ProfitLoss,
				ProfitLossRate:         upbitInvestment.ProfitLossRate,
			})
		}
	} else {
		BadRequest(c,
			fmt.Sprintf("%s: Unknown assetType", assetType))
	}
	c.IndentedJSON(http.StatusOK, investments)
}

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

	router := gin.Default()
	router.GET("/investments", getInvestments)
	router.Run("localhost:61001")
}
