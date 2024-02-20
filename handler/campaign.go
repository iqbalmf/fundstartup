package handler

import (
	"fmt"
	"funding-app/campaign"
	"funding-app/helper"
	"funding-app/users"
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

// GetCampaign api/v1/campaigns/"id"
func (r *campaignHandler) GetCampaign(c *gin.Context) {
	//handler: mapping id dari url ke struct input call service, call formatter
	//service : struct input menangkap id di url, call repo
	//repository: get campaign by id
	var input campaign.GetCampaignDetailInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail campaign", http.StatusUnprocessableEntity, "error", nil)
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

// CREATE CAMPAIGN
// get param from user input struct
// get current user from jwt/handler
// call service, param's input struct (and create slug)
// call repository to save new campaign
func (r *campaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := helper.APIResponse("Failed to Create Campaign", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	currentUser := c.MustGet("currentUser").(users.User)
	input.User = currentUser
	newCampaign, err := r.service.CreateCampaign(input)
	if err != nil {
		response := helper.APIResponse("Failed to Create Campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Success to create new campaign", http.StatusOK, "success", campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusOK, response)
}

// UPDATE CAMPAIGN
// user input data
// handler
// mapping from input to input struct
// service
// repository to save db
func (r *campaignHandler) UpdateCampaign(c *gin.Context) {
	var inputID campaign.GetCampaignDetailInput
	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.APIResponse("Failed to Update Campaign", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	var inputData campaign.CreateCampaignInput
	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		response := helper.APIResponse("Failed to Update Campaign", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	currentUser := c.MustGet("currentUser").(users.User)
	inputData.User = currentUser
	updateCampaign, err := r.service.UpdateCampaign(inputID, inputData)
	if err != nil {
		response := helper.APIResponse(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Success to Update campaign", http.StatusOK, "success", campaign.FormatCampaign(updateCampaign))
	c.JSON(http.StatusOK, response)
}

// Upload Image Campaign
// handler
// get input change to struct input
// save image to folder
// service
// repository :
// save data image to table campaign_images
// change is_primary true to false
func (r *campaignHandler) UploadImage(c *gin.Context) {
	var input campaign.CreateCampaignImageInput
	err := c.ShouldBind(&input)
	if err != nil {
		response := helper.APIResponse(err.Error(), http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	currentUser := c.MustGet("currentUser").(users.User)
	input.User = currentUser
	userId := currentUser.ID
	file, err := c.FormFile("file")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse(err.Error(), http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	path := fmt.Sprintf("campaign_images/%d-%s", userId, file.Filename)
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse(err.Error(), http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	_, err = r.service.SaveCampaignImages(input, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse(err.Error(), http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Image success uploaded", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
	return
}
