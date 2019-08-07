package routes

import (
	"net/http"

	"github.com/airnomadsmitty/quarters/utils"
)

type IndexController struct{}

func MakeIndexController() *IndexController {
	return &IndexController{}
}

func (cont *IndexController) Get(res http.ResponseWriter, req *http.Request, auth *utils.Auth) {
	if auth.IsLoggedIn() {
		http.Redirect(res, req, "/dashboard", http.StatusSeeOther)
	}

	http.Redirect(res, req, "/login", http.StatusSeeOther)
}
