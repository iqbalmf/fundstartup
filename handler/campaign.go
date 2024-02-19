package handler

import (
	"funding-app/campaign"
	"funding-app/helper"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// get parameters handler
// handler to service
// service to decide what repository to call
// repository: GetAll, GetByUserID
// db

type campaignHandler struct {
	service campaign.Service
}

func NewUserCampaign(service campaign.Service) *campaignHandler {
	return &campaignHandler{service: service}
}

// GetCampaigns api/v1/campaigns
func (receiver *campaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))
	campaigns, err := receiver.service.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Error to get Campaigns", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
}
