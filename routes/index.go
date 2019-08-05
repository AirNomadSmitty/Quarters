package routes

import (
	"html/template"
	"net/http"
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
	p1 := &indexData{Title: "TestPage", Body: []byte("This is the index my dudes.")}
	t, err := template.ParseFiles("views/index.html")
	if err != nil {
		panic(err.Error())
	}
	t.Execute(res, p1)
}
