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

// GetCampaigns api/v1/campaigns/?user_id="user_id"
func (r *campaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))
	campaigns, err := r.service.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Error to get Campaigns", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
}

// GetCampaign api/v1/campaigns/1
func (r *campaignHandler) GetCampaign(c *gin.Context) {
	//handler: mapping id dari url ke struct input call service, call formatter
	//service : struct input menangkap id di url, call repo
	//repository: get campaign by id
	var input campaign.GetCampaignDetailInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	campaignDetail, err := r.service.GetCampaignByID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if campaignDetail.UserID == 0 {
		response := helper.APIResponse("Detail campaign not found", http.StatusOK, "success", nil)
		c.JSON(http.StatusOK, response)
		return
	}
	response := helper.APIResponse("Campaign detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}
