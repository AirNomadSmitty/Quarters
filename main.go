package main

import (
	"database/sql"
	"net/http"

	"github.com/airnomadsmitty/quarters/mappers"
	"github.com/airnomadsmitty/quarters/utils"

	"github.com/airnomadsmitty/quarters/routes"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	db, err := sql.Open("mysql", "root:@/quarters")
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
	betMapper := mappers.NewBetMapper(db)

	signup := routes.MakeSignupController(userMapper)
	login := routes.MakeLoginController(userMapper)
	index := routes.MakeIndexController()
	dashboard := routes.MakeDashboardController(userMapper, betMapper)
	logout := routes.MakeLogoutController()
	bet := routes.NewBetController(betMapper)
	createBet := routes.NewCreateBetController(betMapper)

	r.HandleFunc("/", NewAuthenticatedWrapper(index.Get).ServeHTTP).Methods("GET")

	r.HandleFunc("/signup", signup.Get).Methods("GET")
	r.HandleFunc("/signup", signup.Post).Methods("POST")

	r.HandleFunc("/login", login.Get).Methods("GET")
	r.HandleFunc("/login", login.Post).Methods("POST")
	r.HandleFunc("/logout", NewAuthenticatedWrapper(logout.Get).ServeHTTP).Methods("GET")

	r.HandleFunc("/dashboard", NewAuthenticatedWrapper(dashboard.Get).ServeHTTP).Methods("GET")

	r.HandleFunc("/bet/create", NewAuthenticatedWrapper(createBet.Get).ServeHTTP).Methods("GET")
	r.HandleFunc("/bet/create", NewAuthenticatedWrapper(createBet.Post).ServeHTTP).Methods("POST")
	r.HandleFunc("/bet/{id:[0-9]+}", NewAuthenticatedWrapper(bet.Get).ServeHTTP).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

type AuthenticatedHandler func(http.ResponseWriter, *http.Request, *utils.Auth)

type AuthenticatedWrapper struct {
	handler AuthenticatedHandler
}

func (wrapper *AuthenticatedWrapper) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	auth, _ := getAuthFromRequest(req)
	wrapper.handler(res, req, auth)
}

func NewAuthenticatedWrapper(handler AuthenticatedHandler) *AuthenticatedWrapper {
	return &AuthenticatedWrapper{handler}
}

func getAuthFromRequest(req *http.Request) (*utils.Auth, error) {
	c, err := req.Cookie("token")
	noAuth := &utils.Auth{}
	if err != nil {
		return noAuth, err
	}

	tokenString := c.Value

	auth, err := utils.GetAuthFromJWT(tokenString)
	if err != nil {
		return noAuth, err
	}

	return auth, nil
}
