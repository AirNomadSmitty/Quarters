package routes

import (
	"net/http"
	"strconv"

	"github.com/airnomadsmitty/quarters/mappers"
	"github.com/airnomadsmitty/quarters/utils"
)

type DashboardController struct {
	userMapper *mappers.UserMapper
	betMapper  *mappers.BetMapper
}

func MakeDashboardController(userMapper *mappers.UserMapper, betMapper *mappers.BetMapper) *DashboardController {
	return &DashboardController{userMapper, betMapper}
}

func (cont *DashboardController) Get(res http.ResponseWriter, req *http.Request, auth *utils.Auth) {
	if !auth.IsLoggedIn() {
		http.Redirect(res, req, "/login", 301)
	}

	
	cont.betMapper.GetFromBetIdsWithUserId([2], auth.UserID)

	res.Write([]byte("Logged in with UserID " + strconv.FormatInt(auth.UserID, 10)))
}
