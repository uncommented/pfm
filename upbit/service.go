package upbit

import (
	"strconv"
)

type UpbitAccountService struct {
	UnimplementedUpbitAccountServer
}

func (ps *UpbitAccountService) ListInvestments(request *UpbitAccountRequest, stream UpbitAccount_ListInvestmentsServer) error {
	investmentsResponse := RequestBalance()

	for _, investment := range investmentsResponse {
		currency := investment["currency"].(string)

		if currency == "KRW" {
			continue
		}

		marketInfo := RequestMarketInfo(currency)
		currency_fullname := marketInfo["english_name"].(string)

		quantity, _ := strconv.ParseFloat(investment["balance"].(string), 64)
		averagePurchasingPrice, _ := strconv.ParseFloat(investment["avg_buy_price"].(string), 64)
		purchasingAmount := averagePurchasingPrice * quantity

		marketSnapshot := RequestMarketSnapshot(currency)
		currentPrice, _ := marketSnapshot["trade_price"].(float64)
		evaluationAmount := currentPrice * quantity
		profitLoss := evaluationAmount - purchasingAmount
		profitLossRate := profitLoss / purchasingAmount
		item := Investment{
			Currency:               currency,
			CurrencyFullname:       currency_fullname,
			Quantity:               quantity,
			AveragePurchasingPrice: averagePurchasingPrice,
			PurchasingAmount:       purchasingAmount,
			CurrentPrice:           currentPrice,
			EvaluationAmount:       evaluationAmount,
			ProfitLoss:             profitLoss,
			ProfitLossRate:         profitLossRate,
		}

		if err := stream.Send(&item); err != nil {
			return err
		}
	}
	return nil
}
