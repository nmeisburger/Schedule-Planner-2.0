package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type unParsedCourse struct {
	Title       string
	Crn         int
	SubjCode    string
	CrseNum     int
	Subject     string
	Department  string
	Credits     string
	Preq        string
	Coreq       string
	Instructors string
	Times       string
	Dists       string
}

type course struct {
	Title       string
	Crn         int
	SubjCode    string
	CrseNum     int
	Subject     string
	Department  string
	Credits     string
	Preq        string
	Coreq       string
	Instructors string
	Times       []courseTimes
	Dists       string
}

var (
	subjRegex  = regexp.MustCompile("[A-Z]{4}")
	crnRegex   = regexp.MustCompile("[0-9]{5}")
	distRegex  = regexp.MustCompile("[0-9]")
	queryRegex = regexp.MustCompile("\\w{1,30}")
)

func (c *unParsedCourse) convert() (course, error) {
	parsedTimes, err := parseTimes(c.Times)
	if err != nil {
		return course{}, err
	}
	return course{c.Title, c.Crn, c.SubjCode, c.CrseNum,
		c.Subject, c.Department, c.Credits, c.Preq, c.Coreq,
		c.Instructors, parsedTimes, c.Dists}, nil
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query, ok := parseQueryParam(r, "query", queryRegex)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	queryString := searchFilters(r)
	results, err := db.Query(queryString)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer results.Close()
	distances := make([]pair, 0, 4000)
	for results.Next() {
		var c unParsedCourse
		results.Scan(&c.Title, &c.Crn, &c.SubjCode, &c.CrseNum, &c.Subject, &c.Department,
			&c.Credits, &c.Preq, &c.Coreq, &c.Instructors, &c.Times, &c.Dists)
		distance := minLevenshteinDistance(*query, &c)
		distances = append(distances, pair{Distance: distance, Course: &c})
	}
	quickSort(&distances)
	courses := make([]course, 0, 30)
	for i := 0; i < min2(len(distances), 30); i++ {
		parsedCourse, err := (*distances[i].Course).convert()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		courses = append(courses, parsedCourse)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func instructorSearchHandler(w http.ResponseWriter, r *http.Request) {
	query, ok := parseQueryParam(r, "instructor", queryRegex)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	queryString := searchFilters(r)
	results, err := db.Query(queryString)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer results.Close()
	distances := make([]pair, 0, 4000)
	for results.Next() {
		var c unParsedCourse
		results.Scan(&c.Title, &c.Crn, &c.SubjCode, &c.CrseNum, &c.Subject, &c.Department,
			&c.Credits, &c.Preq, &c.Coreq, &c.Instructors, &c.Times, &c.Dists)
		distance := levenshteinDistance(*query, strings.ToLower(c.Instructors))
		distances = append(distances, pair{Distance: distance, Course: &c})
	}
	quickSort(&distances)
	courses := make([]course, 0, 30)
	for i := 0; i < min2(len(distances), 30); i++ {
		parsedCourse, err := (*distances[i].Course).convert()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		courses = append(courses, parsedCourse)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func crnQueryHandler(w http.ResponseWriter, r *http.Request) {
	crn, ok := parseQueryParam(r, "crn", crnRegex)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	q := fmt.Sprintf("SELECT * FROM courses c WHERE c.Crn = %s", *crn)
	result := db.QueryRow(q)
	var c course
	var times string
	result.Scan(&c.Title, &c.Crn, &c.SubjCode, &c.CrseNum, &c.Subject, &c.Department,
		&c.Credits, &c.Preq, &c.Coreq, &c.Instructors, &times, &c.Dists)
	parsedTimes, err := parseTimes(times)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Times = parsedTimes
	var courses = []course{c}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func filterQueryHandler(w http.ResponseWriter, r *http.Request) {
	subj, subjOk := parseQueryParam(r, "subj", subjRegex)
	dist, distOk := parseQueryParam(r, "dist", distRegex)
	if !subjOk && !distOk {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var q string
	if !subjOk && distOk {
		q = fmt.Sprintf("SELECT * FROM courses c WHERE c.Dists = '%s'", *dist)
	} else if subjOk && !distOk {
		q = fmt.Sprintf("SELECT * FROM courses c WHERE c.SubjCode = '%s'", *subj)
	} else {
		q = fmt.Sprintf("SELECT * FROM courses c WHERE c.SubjCode = '%s' AND c.Dists = '%s'", *subj, *dist)

	}
	results, err := db.Query(q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer results.Close()
	courses := make([]course, 0, 100)
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

func parseQueryParam(r *http.Request, key string, re *regexp.Regexp) (*string, bool) {
	val, ok := r.URL.Query()[key]
	if !ok {
		return nil, false
	}
	m := re.FindStringSubmatch(val[0])
	if m == nil || len(m) != 1 {
		return nil, false
	}
	return &m[0], true
}

func searchFilters(r *http.Request) string {
	subj, subjOk := parseQueryParam(r, "subj", subjRegex)
	dist, distOk := parseQueryParam(r, "dist", distRegex)
	if !subjOk && !distOk {
		return "SELECT * FROM courses"
	}
	if !subjOk && distOk {
		return fmt.Sprintf("SELECT * FROM courses c WHERE c.Dists = '%s'", *dist)
	} else if subjOk && !distOk {
		return fmt.Sprintf("SELECT * FROM courses c WHERE c.SubjCode = '%s'", *subj)
	} else {
		return fmt.Sprintf("SELECT * FROM courses c WHERE c.SubjCode = '%s' AND c.Dists = '%s'", *subj, *dist)
	}
}
