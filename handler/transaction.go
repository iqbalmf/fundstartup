package handler

import (
	"funding-app/helper"
	"funding-app/transaction"
	"funding-app/users"
	"github.com/gin-gonic/gin"
	"net/http"
)

// parameter di uri
// getting param, mapping to input struct
// call service, input struct as param
// service with campaignID call repository
// repo find data transaction's campaign
type transactionHandler struct {
	service transaction.Service
}

func NewTransaction(service transaction.Service) *transactionHandler {
	return &transactionHandler{service: service}
}
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
