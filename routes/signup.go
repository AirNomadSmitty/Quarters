package routes

import (
	"net/http"

	"github.com/airnomadsmitty/quarters/mappers"

	"golang.org/x/crypto/bcrypt"
)

type SignupController struct {
	userMapper *mappers.UserMapper
}

func MakeSignupController(userMapper *mappers.UserMapper) *SignupController {
	return &SignupController{userMapper}
}

func (cont *SignupController) Get(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "views/signup.html")
}

func (cont *SignupController) Post(res http.ResponseWriter, req *http.Request) {
	username := req.FormValue("username")
	password := req.FormValue("password")

	user, err := cont.userMapper.GetFromUsername(username)

	switch {
	case err != nil:
		http.Error(res, "Server error.", 500)
		return
	case user == nil:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}

		user, err = cont.userMapper.Create(username, string(hashedPassword))
		if err != nil {
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}

		res.Write([]byte("User created!"))
		return
	default:
		http.Error(res, "User already exists", 500)
	}
}
