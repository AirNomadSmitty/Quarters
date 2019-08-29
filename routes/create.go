package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/airnomadsmitty/quarters/mappers"

	"github.com/airnomadsmitty/quarters/utils"
)

type CreateBetController struct {
	betMapper *mappers.BetMapper
}

func NewCreateBetController(betmapper *mappers.BetMapper) *CreateBetController {
	return &CreateBetController{}
}

func (cont *CreateBetController) Get(res http.ResponseWriter, req *http.Request, auth *utils.Auth) {
	http.ServeFile(res, req, "views/bet/create.html")
}

func (cont *CreateBetController) Post(res http.ResponseWriter, req *http.Request, auth *utils.Auth) {
	position := req.FormValue("position")
	baseValue, err := strconv.ParseFloat(req.FormValue("baseValue"), 64)
	if err != nil {
		panic(err)
	}
	note := req.FormValue("note")
	oddsMultiplier, err := strconv.ParseFloat(req.FormValue("oddsMultiplier"), 64)
	if err != nil {
		panic(err)
	}

	fmt.Println("creating bet...")
	newBet, err := cont.betMapper.Create(baseValue, position, oddsMultiplier, auth.UserID, note)
	if err != nil {
		panic(err)
	}

	http.Redirect(res, req, "/bet/"+strconv.FormatInt(newBet.BetID, 10), http.StatusSeeOther)
}
