package main

import (
	"database/sql"
	"net/http"

	"github.com/airnomadsmitty/quarters/mappers"

	"github.com/airnomadsmitty/quarters/routes"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

type Home struct {
	Title string
	Body  []byte
}

func main() {
	db, err = sql.Open("mysql", "root:@/quarters")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	r := mux.NewRouter()

	userMapper := mappers.MakeUserMapper(db)

	signup := routes.MakeSignupController(userMapper)
	login := routes.MakeLoginController(userMapper)
	index := routes.MakeIndexController()

	r.HandleFunc("/", index.Get).Methods("GET")

	r.HandleFunc("/signup", signup.Get).Methods("GET")
	r.HandleFunc("/signup", signup.Post).Methods("POST")

	r.HandleFunc("/login", login.Get).Methods("GET")
	r.HandleFunc("/login", login.Post).Methods("POST")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
