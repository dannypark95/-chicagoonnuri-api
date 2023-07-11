package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dannypark95/ChicagoOnnuri/models"
	"github.com/dgrijalva/jwt-go"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	
	// Decode the reuqest body to get user's email and password
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
		return
	}

	// Log the request
	fmt.Printf("Login request: username=%s\n", creds.Username)

	// Validate user's email and password
	user, err := models.AuthenticateUser(creds.Username, creds.Password)
	if err != nil {
		http.Error(w, `{"error": "Invalid username or password"}`, http.StatusUnauthorized)
		return
	}

	// Create a new token object
	token := jwt.New(jwt.SigningMethodHS256)

	// Create a map to store our claims
	claims := token.Claims.(jwt.MapClaims)

	// Set token claims
	claims["username"] = user.Username
	claims["admin"] = user.Admin
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Sign the token with our secret
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	// Finally, write the token to the browser window
	w.Write([]byte(tokenString))
}