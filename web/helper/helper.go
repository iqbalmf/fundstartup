package helper

import (
	"funding-app/helper"
	"github.com/gorilla/sessions"
)

var CookieStore = sessions.NewCookieStore([]byte(helper.GoDotEnvVariable("SECRET_KEY_JWT")))
