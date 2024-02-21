package handler

import (
	"funding-app/helper"
	"funding-app/transaction"
	"funding-app/users"
	"github.com/gin-gonic/gin"
	"net/http"
)

type transactionHandler struct {
	service transaction.Service
}

func NewTransaction(service transaction.Service) *transactionHandler {
	return &transactionHandler{service: service}
}

// parameter di uri
// getting param, mapping to input struct
// call service, input struct as param
// service with campaignID call repository
// repo find data transaction's campaign
func (t *transactionHandler) GetTransactionCampaign(c *gin.Context) {
	var input transaction.GetTransactionCampaignInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Error to get Transaction Campaign's", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	currentUser := c.MustGet("currentUser").(users.User)
	input.User = currentUser
	tr, err := t.service.GetTransactionByCampaignID(input)
	if err != nil {
		response := helper.APIResponse("Error to get Transaction Campaign's. "+err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Transaction Campaign's detail", http.StatusOK, "success", transaction.FormatTransactionCampaigns(tr))
	c.JSON(http.StatusOK, response)
}

// GetTeransactionUser
// handler
// get user from jwt/middleware
// service
// repo => get data transaction(preload campaign)
func (t *transactionHandler) GetTransactionUser(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(users.User)
	userId := currentUser.ID
	tr, err := t.service.GetTransactionByUserID(userId)
	if err != nil {
		response := helper.APIResponse("Error to get Transaction User's "+err.Error(), http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	response := helper.APIResponse("Transaction User's detail", http.StatusOK, "success", transaction.FormatTransactions(tr))
	c.JSON(http.StatusOK, response)
}

// MIDTRANS Integrations
// input from user
// handler get input and pass to input struct
// service to create transaction, call system midtrans
// repository to create new transaction
func (t *transactionHandler) CreateTransaction(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(users.User)
	var input transaction.CreateTransactionInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		error := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": error}
		response := helper.APIResponse("Failed to create Transaction", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	input.User = currentUser
	newTransaction, err := t.service.SaveTransaction(input)
	if err != nil {
		response := helper.APIResponse("Failed to create Transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Success to create transaction", http.StatusOK, "error", transaction.FormatPaymentTransaction(newTransaction))
	c.JSON(http.StatusOK, response)
}
