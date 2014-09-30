package routes

import (
	"net/http"

	"code.google.com/p/go-uuid/uuid"
	"github.com/lavab/api/db"
	"github.com/lavab/api/models"
	"github.com/lavab/api/utils"
)

// Login TODO
func Login(w http.ResponseWriter, r *http.Request) {
	user, pass := r.FormValue("username"), r.FormValue("password")
	loginHelper(user, pass, w, r)
}

// Signup TODO
func Signup(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("username")
	pass := r.FormValue("password")
	// regt := r.FormValue("reg_token")

	if _, err := db.User(user); err == nil {
		utils.JSONResponse(w, map[string]interface{}{
			"status":  409,
			"success": false,
			"message": "Username already exists",
		})
		return
	}
	hash, err := utils.BcryptHash(pass)
	if err != nil {
		utils.JSONResponse(w, map[string]interface{}{
			"status":  500,
			"message": "Hashing with bcrypt has failed",
		})
		return
	}
	created := models.User{
		ID:       uuid.New(),
		Name:     user,
		Password: string(hash),
	}
	err = db.CreateUser(created)
	if err != nil {
		utils.JSONResponse(w, map[string]interface{}{
			"status":  500,
			"message": "Couldn't save the data to database",
		})
		return
	}

	loginHelper(user, pass, w, r)
}

// Logout TODO
func Logout(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, map[string]interface{}{
		"status":  404,
		"message": "Hey Dennis, this isn't implemented yet! :D",
	})
}

// Me TODO
func Me(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	// TODO make this check a middleware function
	if token == "" {
		utils.JSONResponse(w, map[string]interface{}{
			"status":  401,
			"message": "Please login to view this resource",
		})
		return
	}
	session, err := db.Session(token)
	if err != nil {
		utils.JSONResponse(w, map[string]interface{}{
			"status":  401,
			"message": "Invalid token",
		})
		return
	}
	user, err := db.User(session.User)
	if err != nil {
		utils.JSONResponse(w, map[string]interface{}{
			"status":  500,
			"message": "Corrupted session or user store",
		})
		return
	}
	utils.JSONResponse(w, map[string]interface{}{
		"status": 200,
		"data":   user,
	})
}

func loginHelper(user, pass string, w http.ResponseWriter, r *http.Request) {
	userData, err := db.User(user)
	if err != nil {
		utils.JSONResponse(w, map[string]interface{}{
			"status":  403,
			"message": "Wrong username of password",
		})
		return
	}
	if !utils.BcryptVerify(userData.Password, pass) {
		utils.JSONResponse(w, map[string]interface{}{
			"status":  403,
			"message": "Wrong username of password",
		})
		return
	}
	token, err := db.CreateSession(user, 72)
	if err != nil {
		utils.JSONResponse(w, map[string]interface{}{
			"status":  500,
			"message": "Unable to create session",
		})
		return
	}
	// For now we're sending the token in plaintext, until I implement JWT
	utils.JSONResponse(w, map[string]interface{}{
		"status":   200,
		"success":  true,
		"token":    token,
		"exp_date": utils.HoursFromNow(72),
	})
}
