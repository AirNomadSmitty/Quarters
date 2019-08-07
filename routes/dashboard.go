package routes

import (
	"net/http"
	"strconv"

	"github.com/airnomadsmitty/quarters/mappers"
	"github.com/airnomadsmitty/quarters/utils"
)

type DashboardController struct {
	userMapper *mappers.UserMapper
}

func MakeDashboardController(userMapper *mappers.UserMapper) *DashboardController {
	return &DashboardController{userMapper}
}

func (cont *DashboardController) Get(res http.ResponseWriter, req *http.Request, auth *utils.Auth) {
	if !auth.IsLoggedIn() {
		http.Redirect(res, req, "/login", 301)
	}

	res.Write([]byte("Logged in with UserID " + strconv.FormatInt(auth.UserID, 10)))
}
