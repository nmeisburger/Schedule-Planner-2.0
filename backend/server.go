package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/handlers"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db, dbConnectionError = sql.Open("mysql", dbAddress)

// var subjectTrie trie

func main() {
	if dbConnectionError != nil {
		panic(dbConnectionError)
	}

	result := db.QueryRow("SELECT COUNT(*) FROM courses")
	var x int
	result.Scan(&x)
	log.Printf("Total Courses: %d", x)

	router := mux.NewRouter()

	router.HandleFunc("/filter", filterQueryHandler).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/crn", crnQueryHandler).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/search", searchHandler).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/instructor", instructorSearchHandler).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/register", registerHandler).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/signin", signinHandler).Methods(http.MethodPost, http.MethodOptions)

	router.HandleFunc("/addcourse", addCourseHandler).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/removecourse", removeCourseHandler).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/workspace", workspaceHandler).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/trie", handleTrie).Methods(http.MethodGet, http.MethodOptions)

	log.Println("Server Launched")
	log.Fatal(http.ListenAndServe(":3002",
		handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Access-Token"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
			handlers.ExposedHeaders([]string{"Access-Token"}),
			handlers.AllowedOrigins([]string{"*"}))(router)))
}
