package handler

import (
	"funding-app/users"
	"github.com/gin-gonic/gin"
	"net/http"
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
	}
	c.HTML(http.StatusOK, "user_index.html", gin.H{"users": users})
}
