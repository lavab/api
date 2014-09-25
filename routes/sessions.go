package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/lavab/api/db"
	"github.com/lavab/api/models"
	"github.com/lavab/api/util"
)

const cost = bcrypt.DefaultCost * 2

type loginResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Expires string `json:"expires"`
	Message string `json:"message"`
}

// Login TODO
func Login(w http.ResponseWriter, r *http.Request) {
	user, pass := r.FormValue("username"), r.FormValue("password")

	login(user, pass, w, r)
}

func login(user, pass string, w http.ResponseWriter, r *http.Request) {
	userData, err := db.GetUser(user)
	if err != nil {
		http.Error(w, "Wrong username or password", 403)
		return
	}
	if userData.Salt == "" {
		http.Error(w, "The user doesn't have a salt", 500)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pass+userData.Salt), cost)
	if err != nil {
		http.Error(w, "Error salting the provided password", 500)
		return
	}
	if string(hash) != userData.Password {
		http.Error(w, "Wrong username or password", 403)
		return
	}

	token, err := db.CreateSession(user, userData.ID, r.UserAgent())
	if err != nil {
		http.Error(w, "Unable to create session", 500)
		return
	}

	// For now we're sending the token in plaintext, until I implement JWT
	res, err := json.Marshal(loginResponse{
		Success: true,
		Token:   token,
		Expires: util.HoursFromNow(80),
		Message: "",
	})
	if err != nil {
		http.Error(w, "Error marshaling the response body", 500)
		return
	}
	fmt.Fprintf(w, string(res))
}

// Signup TODO
func Signup(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("username")
	pass := r.FormValue("password")
	// regt := r.FormValue("reg_token")

	// TODO duplicate code
	if _, err := db.GetUser(user); err == nil {
		res, err := json.Marshal(loginResponse{
			Success: false,
			Token:   "",
			Expires: "",
			Message: "Username already exists",
		})
		if err != nil {
			http.Error(w, "Error marshaling the response body", 500)
			return
		}
		fmt.Fprintf(w, string(res))
		return
	}

	salt, err := util.RandomString(16)
	if err != nil {
		http.Error(w, "Error generating secure random string", 500)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pass+salt), cost)
	if err != nil {
		http.Error(w, "Error salting the provided password", 500)
		return
	}

	created := models.User{
		ID:       uuid.New(),
		Name:     user,
		Password: string(hash),
		Salt:     salt,
	}
	err = db.CreateUser(created)
	if err != nil {
		http.Error(w, "Error saving user to database", 500)
		return
	}

	login(user, pass, w, r)
}

// Logout TODO
func Logout(w http.ResponseWriter, r *http.Request) {

}

// Me TODO
func Me(w http.ResponseWriter, r *http.Request) {

}
