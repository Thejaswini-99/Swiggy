package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

	fmt.Println("CacheHit-- Fetching the token from cache")
	val, err := rdb.Get(ctx, "token:"+loginUser.Username).Result()
	if err == nil {
		fmt.Println("CacheHit- Retiring from the redis")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", val)
		return
	}

	token, err := generateToken(loginUser.Username)
	if err != nil {
		fmt.Fprintf(w, "Error in Generating  the token")
		return
	}
	fmt.Println("Cache Miss-- storing the token in cache")
	rdb.Set(ctx, "token:"+loginUser.Username, token, 24*time.Hour)
	fmt.Println("Stored in Redis Cache")

	fmt.Fprint(w, token)
}

func logout(w http.ResponseWriter, r *http.Request) {
	tokenstr := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	username, err := getUsernameFromToken(tokenstr)
	if err != nil {
		fmt.Fprintf(w, "Invalid token!")
		return
	}
	rdb.Del(ctx, "token:"+username)

	fmt.Fprintf(w, "Logged Out Successfully")

}

func getUsernameFromToken(tokenstr string) (string, error) {
	token, err := validateToken(tokenstr)
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("Invlaid Claims")
	}
	username, ok := claims["username"].(string)
	if !ok {
		return "", fmt.Errorf("Username not found")
	}

	return username, nil
}
