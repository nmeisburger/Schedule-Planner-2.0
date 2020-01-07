package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var secret = []byte("20sdf0")

var users = map[string]string{"u1": "p1", "u2": "p2"}

func signin(w http.ResponseWriter, r *http.Request) {
	username, password := r.FormValue("username"), r.FormValue("password")
	fmt.Println(username, password)
	credentials := credentials{Username: username, Password: password}
	if len(username) == 0 || len(password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	correct, ok := users[credentials.Username]
	if !ok || correct != credentials.Password {
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
		Name:    "accessToken",
		Value:   tokenString,
		Expires: expiration})
}

func thing(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("accessToken")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}
	tokenString := c.Value
	claims := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) { return secret, nil })
	if err != nil || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
	}
	w.Write([]byte(fmt.Sprintf("Access Granted")))
}

// func main() {
// 	http.HandleFunc("/signin", signin)
// 	http.HandleFunc("/thing", thing)

// 	log.Fatal(http.ListenAndServe(":3000", nil))
// }
