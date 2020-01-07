package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type coursesXML struct {
	XMLName xml.Name    `xml:"COURSES"`
	Courses []courseXML `xml:"COURSE"`
}

type courseXML struct {
	XMLName     xml.Name       `xml:"COURSE"`
	Crn         string         `xml:"crn,attr"`
	SubjCode    string         `xml:"subj-code,attr"`
	CrseNum     string         `xml:"crse-numb,attr"`
	Title       string         `xml:"CRSE_TITLE"`
	Subject     string         `xml:"SUBJECT"`
	Department  string         `xml:"DEPARTMENT"`
	Credits     string         `xml:"CREDITS"`
	Preq        string         `xml:"PREQ"`
	Coreq       coreqsXML      `xml:"COREQS"`
	Instructors instructorsXML `xml:"INSTRUCTORS"`
	Times       timesXML       `xml:"TIMES"`
	Dists       distsXML       `xml:"DISTS"`
}

type timesXML struct {
	XMLName xml.Name        `xml:"TIMES"`
	Times   []courseTimeXML `xml:"MEETING"`
}

type courseTimeXML struct {
	XMLName xml.Name `xml:"MEETING"`
	Begin   string   `xml:"begin-time,attr"`
	End     string   `xml:"end-time,attr"`
	Type    string   `xml:"TYPE"`
	Mon     xml.Name `xml:"MON_DAY"`
	Tue     xml.Name `xml:"TUE_DAY"`
	Wed     xml.Name `xml:"WED_DAY"`
	Thu     xml.Name `xml:"THU_DAY"`
	Fri     xml.Name `xml:"FRI_DAY"`
	Sat     xml.Name `xml:"SAT_DAY"`
	Sun     xml.Name `xml:"SUN_DAY"`
}

type instructorsXML struct {
	XMLName xml.Name `xml:"INSTRUCTORS"`
	Names   []string `xml:"NAME"`
}

type coreqsXML struct {
	XMLName xml.Name   `xml:"COREQS"`
	Coreqs  []coreqXML `xml:"COREQ"`
}

type coreqXML struct {
	SubjCode string `xml:"SUBJ,attr"`
	CrseNum  string `xml:"NUMB,attr"`
}

type distsXML struct {
	XMLName      xml.Name `xml:"DISTS"`
	Distrubtions []string `xml:"DIST"`
}

func (c *courseXML) hasPreq() bool {
	if len(c.Preq) > 0 {
		return true
	}
	return false
}

func (c *courseXML) print() {
	fmt.Printf(
		"\nCourse: %v, %v %v\nCRN: %v\nDepartment: %v\nSubject: %v\nInstructors: %v\nCredits: %v\nPrequisites: %v\nCorequisites: %v\nTimes: %v\nDistributions: %v\n\n",
		c.Title, c.SubjCode, c.CrseNum, c.Crn, c.Department, c.Subject, c.Instructors.toString(), c.Credits, c.Preq, c.Coreq.toString(), c.Times.toString(), c.Dists.toString())
}

func (t timesXML) toString() string {
	if len(t.Times) == 0 {
		return ""
	}
	stub := ""
	for _, v := range t.Times {
		if len(v.Begin) == 0 || v.Type != "Class" {
			continue
		}
		time := v.Begin + ";" + v.End + ";"
		if len(v.Mon.Local) > 0 {
			time += "M"
		}
		if len(v.Tue.Local) > 0 {
			time += "T"
		}
		if len(v.Wed.Local) > 0 {
			time += "W"
		}
		if len(v.Thu.Local) > 0 {
			time += "H"
		}
		if len(v.Fri.Local) > 0 {
			time += "F"
		}
		if len(v.Sat.Local) > 0 {
			time += "S"
		}
		if len(v.Sun.Local) > 0 {
			time += "U"
		}
		stub += time + ";"
	}
	if len(stub) == 0 {
		return ""
	}
	return stub[:len(stub)-1]
}

func (i instructorsXML) toString() string {
	if len(i.Names) == 0 {
		return "NA"
	}
	stub := ""
	for _, v := range i.Names {
		stub += v + ";"
	}
	return stub[:len(stub)-1]
}

func (c coreqsXML) toString() string {
	if len(c.Coreqs) == 0 {
		return "NA"
	}
	stub := ""
	for _, v := range c.Coreqs {
		stub += v.SubjCode + " " + v.CrseNum + ";"
	}
	return stub[:len(stub)-1]
}

func (d distsXML) toString() string {
	if len(d.Distrubtions) == 0 {
		return "None"
	}
	stub := ""
	for _, v := range d.Distrubtions {
		switch v {
		case "Distribution Group I":
			stub += "1;"
		case "Distribution Group II":
			stub += "2;"
		case "Distribution Group III":
			stub += "3;"
		}
	}
	return stub[:len(stub)-1]
}

func readFromFile(filename string) (*[]courseXML, error) {
	// File must be opened and saved before parsing (in vs code for some reason)
	xmlFile, err := os.Open(filename)
	defer xmlFile.Close()
	if err != nil {
		log.Println("Read Failed")
		return nil, errors.New("Read Failed")
	}
	bytevalue, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		log.Println("Parse Failed")
		return nil, errors.New("Parse Failed")
	}
	var courses coursesXML
	xml.Unmarshal(bytevalue, &courses)
	log.Printf("Courses Parsed: %v\n", len(courses.Courses))
	// maxTitle, maxSubject, maxDepartment, maxPreq, maxCoreq, maxInstructors, maxTimes := 0, 0, 0, 0, 0, 0, 0
	// for _, course := range courses.Courses {
	// 	if len(course.Title) > maxTitle {
	// 		maxTitle = len(course.Title)
	// 	}
	// 	if len(course.Subject) > maxSubject {
	// 		maxSubject = len(course.Subject)
	// 	}
	// 	if len(course.Department) > maxDepartment {
	// 		maxDepartment = len(course.Department)
	// 	}
	// 	if len(course.Preq) > maxPreq {
	// 		maxPreq = len(course.Preq)
	// 	}
	// 	if len(course.Coreq.toString()) > maxCoreq {
	// 		maxCoreq = len(course.Coreq.toString())
	// 	}
	// 	if len(course.Instructors.toString()) > maxInstructors {
	// 		maxInstructors = len(course.Instructors.toString())
	// 	}
	// 	if len(course.Times.toString()) > maxTimes {
	// 		maxTimes = len(course.Times.toString())
	// 	}
	// }
	// fmt.Printf("Max Title: %v\n", maxTitle)
	// fmt.Printf("Max Subject: %v\n", maxSubject)
	// fmt.Printf("Max Department: %v\n", maxDepartment)
	// fmt.Printf("Max Preq: %v\n", maxPreq)
	// fmt.Printf("Max Coreq: %v\n", maxCoreq)
	// fmt.Printf("Max Instructors: %v\n", maxInstructors)
	// fmt.Printf("Max Times: %v\n", maxTimes)

	return &courses.Courses, nil
}

func getAndRead(term string, year int) (*[]courseXML, error) {
	var s string
	switch term {
	case "fall":
		s = strconv.Itoa(year+1) + "10"
	case "spring":
		s = strconv.Itoa(year) + "20"
	case "summer":
		s = strconv.Itoa(year) + "30"
	}
	response, err := http.Get("https://courses.rice.edu/courses/!swkscat.cat?format=XML&p_action=COURSE&p_term=" + s)

	// https://b9stu-reg-esther202010.rice.edu/StudentRegistrationSsb/ssb/searchResults/searchResults?txt_subject=COMP&txt_term=202030

	if err != nil {
		log.Println("Fetch Failed")
		return nil, errors.New("Fetch Failed")
	}
	defer response.Body.Close()
	bytevalue, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Parse Failed")
		return nil, errors.New("Parse Failed")
	}
	var fetchedCourses coursesXML
	xml.Unmarshal(bytevalue, &fetchedCourses)
	log.Printf("Courses Parsed: %v\n", len(fetchedCourses.Courses))
	return &fetchedCourses.Courses, nil
}

// func main() {
// 	readFromFile("../RiceCourseAPI/courses.xml")
// 	// getAndRead("fall", 2019)
// }
