package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func addCourseHandler(w http.ResponseWriter, r *http.Request) {
	username, err := verifyTokenGetUsername(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	crn := r.FormValue("crn")
	if len(crn) != 5 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	insertion := fmt.Sprintf("INSERT INTO workspaces VALUES('%s', %s)", username, crn)
	_, err = db.Query(insertion)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func removeCourseHandler(w http.ResponseWriter, r *http.Request) {
	username, err := verifyTokenGetUsername(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	crn := r.FormValue("crn")
	if len(username) == 0 || len(crn) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	deletion := fmt.Sprintf("DELETE FROM workspaces w WHERE w.Username = '%s' AND w.Crn = %s", username, crn)
	_, err = db.Query(deletion)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusAccepted)
}

func workspaceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	username, err := verifyTokenGetUsername(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	q := fmt.Sprintf("SELECT * FROM courses c WHERE EXISTS ( SELECT * FROM workspaces w WHERE w.Crn = c.Crn AND w.Username = '%s')", username)
	results, err := db.Query(q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer results.Close()
	courses := make([]course, 0, 10)
	for results.Next() {
		var c course
		var times string
		results.Scan(&c.Title, &c.Crn, &c.SubjCode, &c.CrseNum, &c.Subject, &c.Department,
			&c.Credits, &c.Preq, &c.Coreq, &c.Instructors, &times, &c.Dists)
		parsedTimes, err := parseTimes(times)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		c.Times = parsedTimes
		courses = append(courses, c)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}
