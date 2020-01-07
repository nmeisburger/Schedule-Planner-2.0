package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
)

var secret = []byte("aE832f2z0C3")

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type errorResponse struct {
	Error string `json:"error"`
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	username, password := r.FormValue("username"), r.FormValue("password")
	if len(username) == 0 || len(password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	q := fmt.Sprintf("SELECT COUNT(*) FROM users c WHERE c.Username = '%s'", username)
	count := db.QueryRow(q)
	if count == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var x int
	count.Scan(&x)
	if x != 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	insertion := fmt.Sprintf("INSERT INTO users VALUES('%s', '%s')", username, string(hashedPassword))
	_, err = db.Query(insertion)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	username, password := r.FormValue("username"), r.FormValue("password")
	if len(username) == 0 || len(password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	credentials := credentials{Username: username, Password: password}
	q := fmt.Sprintf("SELECT u.Password FROM users u WHERE u.Username = '%s'", username)
	login := db.QueryRow(q)
	if login == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var hashedPassword string
	login.Scan(&hashedPassword)
	valid := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if valid != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	expiration := time.Now().Add(5 * time.Minute)
	claims := &claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration.Unix()}}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "AccessToken",
		Value:   tokenString,
		Expires: expiration})
}

func verifyToken(f func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("AccessToken")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tokenString := c.Value
		claims := &claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) { return secret, nil })
		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		f(w, r)
	}
}

func verifyTokenGetUsername(r *http.Request) (string, error) {
	c, err := r.Cookie("AccessToken")
	if err != nil {
		return "", err
	}
	tokenString := c.Value
	claims := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) { return secret, nil })
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("Invalid Token")
	}
	return claims.Username, nil
}

func sendError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(errorResponse{message})
}

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("approved"))
}
