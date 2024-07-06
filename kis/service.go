package kis

import (
	context "context"

	"github.com/uncommented/pfm/utils"
)

type KISAccountService struct {
	UnimplementedKISAccountServer
}

func (ps *KISAccountService) ListInvestments(request *KISAccountRequest, stream KISAccount_ListInvestmentsServer) error {
	accountNumber := request.AccountNumber
	currency := request.Currency
	marketCode := request.MarketCode
	jsonRes := RequestBalance(accountNumber, currency, marketCode)

	investmentsResponse := utils.UnmarshalToList(jsonRes, "output1")

	for _, _investment := range investmentsResponse {
		if investment, ok := _investment.(map[string]interface{}); ok {
			securityCode := utils.UnmarshalToString(investment, "ovrs_pdno")
			securityFullname := utils.UnmarshalToString(investment, "ovrs_item_name")
			itemType := utils.UnmarshalToInt(investment, "prdt_type_cd")
			quantity := utils.UnmarshalToInt(investment, "ord_psbl_qty")
			averagePurchasingPrice := utils.UnmarshalToFloat(investment, "pchs_avg_pric")
			purchasingAmount := utils.UnmarshalToFloat(investment, "frcr_pchs_amt1")
			currentPrice := utils.UnmarshalToFloat(investment, "now_pric2")
			evaluationAmount := utils.UnmarshalToFloat(investment, "ovrs_stck_evlu_amt")
			profitLoss := utils.UnmarshalToFloat(investment, "frcr_evlu_pfls_amt")
			profitLossRate := utils.UnmarshalToFloat(investment, "evlu_pfls_rt")
			item := Investment{
				SecurityCode:           securityCode,
				SecurityFullname:       securityFullname,
				SecurityTypeCode:       itemType,
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
	}
	return nil
}

func (ps *KISAccountService) GetPerformance(ctx context.Context, request *KISAccountRequest) (*Performance, error) {
	accountNumber := request.AccountNumber
	currency := request.Currency
	marketCode := request.MarketCode
	jsonRes := RequestBalance(accountNumber, currency, marketCode)

	performanceResponse := utils.UnmarshalToMap(jsonRes, "output2")
	purchasingAmount := utils.UnmarshalToFloat(performanceResponse, "frcr_pchs_amt1")
	evaluationAmount := utils.UnmarshalToFloat(performanceResponse, "tot_evlu_pfls_amt")
	profitLoss := utils.UnmarshalToFloat(performanceResponse, "ovrs_tot_pfls")
	profitLossRate := utils.UnmarshalToFloat(performanceResponse, "tot_pftrt")
	performance := Performance{
		PurchasingAmount: purchasingAmount,
		EvaluationAmount: evaluationAmount,
		ProfitLoss:       profitLoss,
		ProfitLossRate:   profitLossRate,
	}
	return &performance, nil
}
