package main

// Powershell Admin:
// $ net start MySql80
// $ net stop MySql80

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

const coursesTable = `CREATE TABLE courses (
	Title VARCHAR(30),
	Crn INTEGER,
	SubjCode VARCHAR(4),
	CrseNum INTEGER,
	Subject VARCHAR(30),
	Department VARCHAR(30),
	Credits VARCHAR(10),
	Preq VARCHAR(180),
	Coreq VARCHAR(16),
	Instructors VARCHAR(120),
	Times VARCHAR(140),
	Dists VARCHAR(5),
	PRIMARY KEY (Crn))`

const usersTable = `CREATE TABLE users (
	Username VARCHAR(30),
	Password VARCHAR(65),
	PRIMARY KEY (Username))`

const workspacesTable = `CREATE TABLE workspaces(
	Username VARCHAR(30),
	Crn INTEGER,
	PRIMARY KEY (Username, Crn))`

func createTable(db *sql.DB, format string, name string) error {
	table, err := db.Query(format)
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("Create Table Failed: %s", name)
	}
	defer table.Close()
	return nil
}

func dropTable(db *sql.DB, name string) error {
	q := fmt.Sprintf("DROP TABLE %s", name)
	result, err := db.Query(q)
	if err != nil {
		return fmt.Errorf("Drop Table Failed: %s", name)
	}
	defer result.Close()
	return nil
}

func insertCourses(courses *[]courseXML, db *sql.DB) error {
	statement, err := db.Prepare("INSERT INTO courses VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return errors.New("Statement Preparation Failed")
	}
	defer statement.Close()
	for _, c := range *courses {
		crn, e1 := strconv.ParseInt(c.Crn, 10, 32)
		crseNum, e2 := strconv.ParseInt(c.CrseNum, 10, 32)
		if e1 != nil || e2 != nil {
			return errors.New("Invalid CRN or Course Number")
		}
		_, err := statement.Exec(
			c.Title, crn, c.SubjCode, crseNum,
			c.Subject, c.Department, c.Credits,
			c.Preq, c.Coreq.toString(),
			c.Instructors.toString(), c.Times.toString(),
			c.Dists.toString())
		if err != nil {
			c.print()
			return err
		}
	}
	return nil
}

// func main() {
// 	db, err := sql.Open("mysql", dbAddress)

// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}

// 	fmt.Println("Connection Successful")

// 	if err = dropCoursesTable(db); err != nil {
// 		fmt.Println(err.Error())
// 	}

// 	if err = initCoursesTable(db); err != nil {
// 		fmt.Println(err.Error())
// 	}

// 	courses, err := readFromFile("../RiceCourseAPI/courses.xml")
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}

// 	fmt.Println("Parse Successful")

// 	if err = insertCourses(courses, db); err != nil {
// 		fmt.Println(err.Error())
// 	}

// 	fmt.Println("Insertion Successful")

// 	defer db.Close()
// }
