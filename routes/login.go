package routes

import (
	"net/http"

	"github.com/airnomadsmitty/quarters/mappers"
	"github.com/airnomadsmitty/quarters/utils"
)

type LoginController struct {
	userMapper *mappers.UserMapper
}

func MakeLoginController(userMapper *mappers.UserMapper) *LoginController {
	return &LoginController{userMapper}
}

func (cont *LoginController) Get(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "views/login.html")
}

func (cont *LoginController) Post(res http.ResponseWriter, req *http.Request) {
	username := req.FormValue("username")
	password := req.FormValue("password")

	user, err := cont.userMapper.GetFromUsername(username)

	if err != nil {
		http.Redirect(res, req, "/login", 301)
		return
	}

	if !user.ValidatePassword(password) {
		http.Redirect(res, req, "/login", 301)
		return
	}

	auth := &utils.Auth{UserID: user.UserID}
	jwt, expirationTime, err := auth.MakeJWT()

	if err != nil {
		http.Redirect(res, req, "/login", 301)
	}

	http.SetCookie(res, &http.Cookie{
		Name:    "token",
		Value:   *jwt,
		Expires: *expirationTime,
	})

	http.Redirect(res, req, "/", 302)
}
