package routes

import (
	"net/http"
	"text/template"

	"github.com/airnomadsmitty/quarters/mappers"
	"github.com/airnomadsmitty/quarters/models"
	"github.com/airnomadsmitty/quarters/utils"
)

type DashboardController struct {
	userMapper *mappers.UserMapper
	betMapper  *mappers.BetMapper
}

type dashboardData struct {
	User *models.User
	Bets map[int64]models.Bet
}

func MakeDashboardController(userMapper *mappers.UserMapper, betMapper *mappers.BetMapper) *DashboardController {
	return &DashboardController{userMapper, betMapper}
}

func (cont *DashboardController) Get(res http.ResponseWriter, req *http.Request, auth *utils.Auth) {
	if !auth.IsLoggedIn() {
		http.Redirect(res, req, "/login", 301)
	}
	var err error
	data := &dashboardData{}
	data.User, err = cont.userMapper.GetFromUserID(auth.UserID)
	if err != nil {
		panic(err.Error())
	}
	data.Bets, err = cont.betMapper.GetFromUserId(data.User.UserID)
	if err != nil {
		panic(err.Error())
	}

	t, err := template.ParseFiles("views/dashboard.html")
	if err != nil {
		panic(err.Error())
	}
	t.Execute(res, data)
}
