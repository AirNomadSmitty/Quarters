package routes

import (
	"net/http"
	"strconv"
	"text/template"

	"github.com/airnomadsmitty/quarters/models"

	"github.com/airnomadsmitty/quarters/mappers"
	"github.com/airnomadsmitty/quarters/utils"
	"github.com/gorilla/mux"
)

type BetController struct {
	betMapper *mappers.BetMapper
}

type betIndexData struct {
	Bet       *models.Bet
	BetString string
}

func NewBetController(betMapper *mappers.BetMapper) *BetController {
	return &BetController{betMapper}
}

func (cont *BetController) Get(res http.ResponseWriter, req *http.Request, auth *utils.Auth) {
	vars := mux.Vars(req)
	data := &betIndexData{}
	betID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		panic(err)
	}
	data.Bet, err = cont.betMapper.GetFromBetIdWithUserId(betID, auth.UserID)
	if err != nil {
		panic(err)
	}
	data.BetString = utils.PrettyPrint(data.Bet)

	t, err := template.ParseFiles("views/bet/index.html")
	if err != nil {
		panic(err.Error())
	}
	t.Execute(res, data)
}
