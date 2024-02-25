package handler

import (
	"fmt"
	"funding-app/users"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type userHandler struct {
	userService users.Service
}

func NewUserHandler(service users.Service) *userHandler {
	return &userHandler{
		service,
	}
}

func (h *userHandler) Index(c *gin.Context) {
	users, err := h.userService.GetAllUser()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.HTML(http.StatusOK, "user_index.html", gin.H{"users": users})
}

func (h *userHandler) NewUser(c *gin.Context) {
	c.HTML(http.StatusOK, "user_new.html", nil)
}

func (h *userHandler) CreateUser(c *gin.Context) {
	var input users.FormCreateUserInput
	err := c.ShouldBind(&input)
	if err != nil {
		input.Error = err
		c.HTML(http.StatusOK, "user_new.html", input)
		return
	}

	registerInput := users.RegisterUserInput{
		Name:       input.Name,
		Email:      input.Email,
		Occupation: input.Occupation,
		Password:   input.Password,
	}
	_, err = h.userService.RegisterUser(registerInput)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.Redirect(http.StatusFound, "/users")
}

func (h *userHandler) GetUserById(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)
	user, err := h.userService.GetUserById(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	input := users.FormUpdateUserInput{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Occupation: user.Occupation,
	}
	c.HTML(http.StatusOK, "user_edit.html", input)
}

func (h *userHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)
	var input users.FormUpdateUserInput
	err := c.ShouldBind(&input)
	if err != nil {
		input.Error = err
		c.HTML(http.StatusOK, "user_edit.html", input)
		return
	}
	input.ID = id
	_, err = h.userService.UpdateUser(input)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.Redirect(http.StatusFound, "/users")
}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)
	c.HTML(http.StatusOK, "user_avatar.html", gin.H{"ID": id})
}

func (h *userHandler) CreatAvatar(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	file, err := c.FormFile("avatar")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	userId := id
	path := fmt.Sprintf("avatar_images/%d-%s", userId, file.Filename)
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	_, err = h.userService.SaveAvatar(userId, path)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.Redirect(http.StatusFound, "/users")
}
