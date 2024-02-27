package handler

import (
	"funding-app/campaign"
	"funding-app/users"
	"github.com/gin-gonic/gin"
	"net/http"
)

type campaignHandler struct {
	campaignSevice campaign.Service
	userService    users.Service
}

func NewCampaignHandler(service campaign.Service, userService users.Service) *campaignHandler {
	return &campaignHandler{service, userService}
}

func (h *campaignHandler) Index(c *gin.Context) {
	campaings, err := h.campaignSevice.GetCampaigns(0)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.HTML(http.StatusOK, "campaign_index.html", gin.H{"campaigns": campaings})
}

func (h *campaignHandler) NewCampaign(c *gin.Context) {
	users, err := h.userService.GetAllUser()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	input := campaign.FormCreateCampaignInput{
		Users: users,
	}
	c.HTML(http.StatusOK, "campaign_new.html", input)
}
