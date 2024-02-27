package handler

import (
	"fmt"
	"funding-app/campaign"
	"funding-app/users"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

func (h *campaignHandler) CreateNewCampaign(c *gin.Context) {
	input := campaign.FormCreateCampaignInput{}
	err := c.ShouldBind(&input)
	if err != nil {
		users, e := h.userService.GetAllUser()
		if e != nil {
			c.HTML(http.StatusInternalServerError, "error.html", nil)
			return
		}
		input.Users = users
		input.Error = err
		c.HTML(http.StatusOK, "campaign_new.html", input)
		return
	}
	user, err := h.userService.GetUserById(input.UserID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	createCampaignInput := campaign.CreateCampaignInput{
		Name:             input.Name,
		ShortDescription: input.ShortDescription,
		Description:      input.Description,
		GoalAmount:       input.GoalAmount,
		Perks:            input.Perks,
		User:             user,
	}
	_, err = h.campaignSevice.CreateCampaign(createCampaignInput)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.Redirect(http.StatusFound, "/campaigns")
}

func (h *campaignHandler) CampaignNewImage(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)
	c.HTML(http.StatusOK, "campaign_image.html", gin.H{"ID": id})
}

func (h *campaignHandler) CreateCampaignNewImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	existCampaign, err := h.campaignSevice.GetCampaignByID(campaign.GetCampaignDetailInput{ID: id})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	userID := existCampaign.UserID
	fmt.Print("user ", userID)
	path := fmt.Sprintf("campaign_images/%d-%s", userID, file.Filename)
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	userCampaign, err := h.userService.GetUserById(userID)
	if err != nil {
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", nil)
			return
		}
	}
	createCampaignImageInput := campaign.CreateCampaignImageInput{
		CampaignID: id,
		IsPrimary:  true,
		User:       userCampaign,
	}
	_, err = h.campaignSevice.SaveCampaignImages(createCampaignImageInput, path)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.Redirect(http.StatusFound, "/campaigns")
}

func (h *campaignHandler) GetCampaignById(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)
	input := campaign.GetCampaignDetailInput{
		ID: id,
	}
	cm, err := h.campaignSevice.GetCampaignByID(input)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	campaignExisting := campaign.FormUpdateCampaignInput{
		ID:               id,
		Name:             cm.Name,
		ShortDescription: cm.ShortDescription,
		Description:      cm.Description,
		GoalAmount:       cm.GoalAmount,
		Perks:            cm.Perks,
	}
	c.HTML(http.StatusOK, "campaign_edit.html", campaignExisting)
}

func (h *campaignHandler) UpdateCampaign(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)
	input := campaign.FormUpdateCampaignInput{
		ID: id,
	}
	err := c.ShouldBind(&input)
	if err != nil {
		input.Error = err
		input.ID = id
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	existCampaign, err := h.campaignSevice.GetCampaignByID(campaign.GetCampaignDetailInput{ID: id})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	userID := existCampaign.UserID
	userCampaign, err := h.userService.GetUserById(userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	updateInput := campaign.CreateCampaignInput{
		Name:             input.Name,
		ShortDescription: input.ShortDescription,
		Description:      input.Description,
		Perks:            input.Perks,
		GoalAmount:       input.GoalAmount,
		User:             userCampaign,
	}
	_, err = h.campaignSevice.UpdateCampaign(campaign.GetCampaignDetailInput{ID: id}, updateInput)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.Redirect(http.StatusFound, "/campaigns")
}

func (h *campaignHandler) ShowCampaignDetail(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)
	existCampaign, err := h.campaignSevice.GetCampaignByID(campaign.GetCampaignDetailInput{ID: id})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.HTML(http.StatusOK, "campaign_show.html", existCampaign)
}
