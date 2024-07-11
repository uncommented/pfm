package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Investment struct {
	ID                     string  `json:"id"`
	Name                   string  `json:"name"`
	Quantity               float64 `json:"quantity"`
	AveragePurchasingPrice float64 `json:"average_purchasing_price"`
	PurchasingAmount       float64 `json:"purchasing_amount"`
	CurrentPrice           float64 `json:"current_price"`
	EvaluationAmount       float64 `json:"evaluation_amount"`
	ProfitLoss             float64 `json:"profit_loss"`
	ProfitLossRate         float64 `json:"profit_loss_rate"`
}

func InternalServerError(c *gin.Context, message string) {
	c.IndentedJSON(http.StatusInternalServerError, gin.H{
		"message": message,
	})
}

func BadRequest(c *gin.Context, message string) {
	c.IndentedJSON(http.StatusInternalServerError, gin.H{
		"message": message,
	})
}
