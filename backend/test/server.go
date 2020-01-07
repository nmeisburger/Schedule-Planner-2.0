package main

import (
	"encoding/json"
	"net/http"
)

type classes struct {
	Title   string
	Cources []class
}

type class struct {
	Name       string
	Crn        int
	Department string
}

func coursesHandler(w http.ResponseWriter, r *http.Request) {
	courses := classes{"Fall Semester 2019",
		[]class{
			class{"Linear Algebra", 355, "MATH"},
			class{"Algorithms", 382, "COMP"}}}
	// js, err := json.Marshal(courses)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	// w.Write(js)
	json.NewEncoder(w).Encode(courses)
}

// func main() {
// 	http.HandleFunc("/courses", coursesHandler)
// 	fmt.Println("Server Launched")
// 	log.Fatal(http.ListenAndServe(":3000", nil))
// }
