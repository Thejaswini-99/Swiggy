package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwt_secret = []byte("dummy_password")

type user struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func generateToken(username string) (string, error) {

	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwt_secret)

}

func validateToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwt_secret, nil
	})

}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Please Enter the valid status Method")
		return
	}
	var newBody user

	err := json.NewDecoder(r.Body).Decode(&newBody)
	if err != nil {
		fmt.Fprintf(w, "Unable to decode the body")
	}
	fmt.Println(newBody.Username)
	fmt.Println(newBody.Password)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newBody.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Fprintf(w, "Error hashing password")
		return
	}

	fmt.Println(hashedPassword)

	_, err = db.Exec("Insert into users(username,password) values(?,?)", newBody.Username, hashedPassword)
	if err != nil {
		fmt.Fprintf(w, "Error Craeting the User")
		return
	}
	fmt.Fprintf(w, "User registered successfully!")
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Please Enter the valid status Method")
		return
	}

	var loginUser user
	err := json.NewDecoder(r.Body).Decode(&loginUser)
	if err != nil {
		fmt.Fprintf(w, "Unable to decode the body")
	}
	fmt.Println(loginUser.Username)
	fmt.Println(loginUser.Password)

	var dbUser user
	err = db.QueryRow("Select username,password from users where username=?", loginUser.Username).Scan(&dbUser.Username, &dbUser.Password)
	if err != nil {
		fmt.Fprintf(w, "Error Fetching the User")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(loginUser.Password))
	if err != nil {
		fmt.Fprintf(w, "Error in comparing the pasword")
		return
	}

	token, err := generateToken(loginUser.Username)
	if err != nil {
		fmt.Fprintf(w, "Error in Generating  the token")
		return
	}

	fmt.Fprintf(w, token)

	// 1. decode request body → get username, password
	// 2. fetch user from DB by username
	// 3. compare password
	// 4. generate JWT token
	// 5. return token
}
