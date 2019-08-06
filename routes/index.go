package routes

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/airnomadsmitty/quarters/utils"
)

type IndexController struct{}

type indexData struct {
	Title string
	Body  []byte
}

func MakeIndexController() *IndexController {
	return &IndexController{}
}

func (cont *IndexController) Get(res http.ResponseWriter, req *http.Request) {
	page := &indexData{}
	c, err := req.Cookie("token")
	if err != nil {
		panic(err.Error())
	}

	tokenString := c.Value

	auth, err := utils.GetAuthFromJWT(tokenString)
	if err != nil {
		panic(err.Error())
	}
	page.Title = "Logged In!"
	page.Body = []byte("Welcome " + strconv.FormatInt(auth.UserID, 10))
	t, err := template.ParseFiles("views/index.html")
	if err != nil {
		panic(err.Error())
	}
	t.Execute(res, page)
}
