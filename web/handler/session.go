package handler

import (
	"errors"
	"funding-app/users"
	webHelper "funding-app/web/helper"
	"github.com/gin-gonic/gin"
	"net/http"
)

type sessionhandler struct {
	userService users.Service
}

func NewSessionHandler(userService users.Service) *sessionhandler {
	return &sessionhandler{userService: userService}
}

func (s *sessionhandler) New(c *gin.Context) {
	c.HTML(http.StatusOK, "session_new.html", nil)
}

func (s *sessionhandler) CreateSession(c *gin.Context) {
	var input users.LoginInput
	err := c.ShouldBind(&input)
	if err != nil {
		input.Error = err
		c.HTML(http.StatusForbidden, "session_new.html", input)
		return
	}
	user, err := s.userService.LoginUser(input)
	if err != nil {
		input.Error = err
		c.HTML(http.StatusForbidden, "session_new.html", input)
		return
	} else if user.Role != "admin" {
		input.Error = errors.New("you're not ADMIN!")
		c.HTML(http.StatusForbidden, "session_new.html", input)
		return
	}

	session, _ := webHelper.CookieStore.Get(c.Request, "login-admin-session")
	session.Values["userID"] = user.ID
	session.Values["nameUser"] = user.Name
	session.Values["emailUser"] = user.Email
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		c.Redirect(http.StatusFound, "error.html")
		return
	}

	c.Redirect(http.StatusFound, "/users")
}

func (s *sessionhandler) Destroy(c *gin.Context) {
	session, err := webHelper.CookieStore.Get(c.Request, "login-admin-session")
	if err != nil {
		c.Redirect(http.StatusFound, "error.html")
		return
	}
	session.Options.MaxAge = -1
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		c.Redirect(http.StatusFound, "error.html")
		return
	}
	c.Redirect(http.StatusFound, "/login")
}
