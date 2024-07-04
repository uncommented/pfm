package kis

import (
	context "context"
	"strconv"
)

type KISAccountService struct {
	UnimplementedKISAccountServer
}

func (ps *KISAccountService) ListInvestments(request *KISAccountRequest, stream KISAccount_ListInvestmentsServer) error {
	accountNumber := request.AccountNumber
	currency := request.Currency
	marketCode := request.MarketCode
	jsonRes := RequestBalance(accountNumber, currency, marketCode)

	investmentsResponse := jsonRes["output1"].([]interface{})

	for _, _investment := range investmentsResponse {
		investment := _investment.(map[string]interface{})
		itemType, _ := strconv.ParseInt(investment["prdt_type_cd"].(string), 10, 64)
		quantity, _ := strconv.ParseInt(investment["ord_psbl_qty"].(string), 10, 64)
		averagePurchasingPrice, _ := strconv.ParseFloat(investment["pchs_avg_pric"].(string), 64)
		purchasingAmount, _ := strconv.ParseFloat(investment["frcr_pchs_amt1"].(string), 64)
		currentPrice, _ := strconv.ParseFloat(investment["now_pric2"].(string), 64)
		evaluationAmount, _ := strconv.ParseFloat(investment["ovrs_stck_evlu_amt"].(string), 64)
		profitLoss, _ := strconv.ParseFloat(investment["frcr_evlu_pfls_amt"].(string), 64)
		profitLossRate, _ := strconv.ParseFloat(investment["evlu_pfls_rt"].(string), 64)
		item := Investment{
			SecurityCode:           investment["ovrs_pdno"].(string),
			SecurityFullname:       investment["ovrs_item_name"].(string),
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
	return nil
}

func (ps *KISAccountService) GetPerformance(ctx context.Context, request *KISAccountRequest) (*Performance, error) {
	accountNumber := request.AccountNumber
	currency := request.Currency
	marketCode := request.MarketCode
	jsonRes := RequestBalance(accountNumber, currency, marketCode)
	performanceResponse := jsonRes["output2"].(map[string]interface{})
	purchasingAmount, _ := strconv.ParseFloat(performanceResponse["frcr_pchs_amt1"].(string), 64)
	evaluationAmount, _ := strconv.ParseFloat(performanceResponse["tot_evlu_pfls_amt"].(string), 64)
	profitLoss, _ := strconv.ParseFloat(performanceResponse["ovrs_tot_pfls"].(string), 64)
	profitLossRate, _ := strconv.ParseFloat(performanceResponse["tot_pftrt"].(string), 64)
	performance := Performance{
		PurchasingAmount: purchasingAmount,
		EvaluationAmount: evaluationAmount,
		ProfitLoss:       profitLoss,
		ProfitLossRate:   profitLossRate,
	}
	return &performance, nil
}
