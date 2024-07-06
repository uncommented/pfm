package upbit

import (
	"github.com/uncommented/pfm/utils"
)

type UpbitAccountService struct {
	UnimplementedUpbitAccountServer
}

func (ps *UpbitAccountService) ListInvestments(request *UpbitAccountRequest, stream UpbitAccount_ListInvestmentsServer) error {
	investmentsResponse := RequestBalance()

	for _, investment := range investmentsResponse {
		currency := utils.UnmarshalToString(investment, "currency")
		if currency == "KRW" {
			continue
		}

		marketInfo := RequestMarketInfo(currency)
		currency_fullname := utils.UnmarshalToString(marketInfo, "english_name")

		quantity := utils.UnmarshalToFloat(investment, "balance")
		averagePurchasingPrice := utils.UnmarshalToFloat(investment, "avg_buy_price")
		purchasingAmount := averagePurchasingPrice * quantity

		marketSnapshot := RequestMarketSnapshot(currency)
		currentPrice := utils.UnmarshalToFloat(marketSnapshot, "trade_price")
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
