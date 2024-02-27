package handler

import (
	"funding-app/transaction"
	"github.com/gin-gonic/gin"
	"net/http"
)

type transactionHandler struct {
	transactionService transaction.Service
}

func NewTransactionHandler(service transaction.Service) *transactionHandler {
	return &transactionHandler{service}
}

func (t *transactionHandler) Index(c *gin.Context) {
	tran, err := t.transactionService.GetAllTransaction()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.HTML(http.StatusOK, "transaction_index.html", gin.H{"transaction": tran})
}
