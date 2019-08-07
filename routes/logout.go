package routes

import (
	"net/http"

	"github.com/airnomadsmitty/quarters/utils"
)

type LogoutController struct{}

func MakeLogoutController() *LogoutController {
	return &LogoutController{}
}

func (cont *LogoutController) Get(res http.ResponseWriter, req *http.Request, auth *utils.Auth) {
	http.SetCookie(res, &http.Cookie{
		Name:   "token",
		Value:  "",
		MaxAge: 0,
	})

	auth.UserID = 0
	http.Redirect(res, req, "/login", 301)
}
