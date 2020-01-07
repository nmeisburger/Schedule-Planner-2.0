package main

// Powershell Admin:
// $ net start MySql80
// $ net stop MySql80

import (
	_ "github.com/go-sql-driver/mysql"
)

type item struct {
	name       string
	department string
	crn        int
}

// func main() {
// 	db, err := sql.Open("mysql", "nmeisburger:salohcin@tcp(127.0.0.1:3306)/courses")

// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	fmt.Println("Connection Successful")

// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	defer db.Close()

// 	// insert1, err1 := db.Query("INSERT INTO courselist VALUES('Linear Algebra', 'MATH', 20932)")
// 	// insert2, err2 := db.Query("INSERT INTO courselist VALUES('Algorithms', 'COMP', 10876)")

// 	// if err1 != nil {
// 	// 	panic(err1.Error())
// 	// }
// 	// defer insert1.Close()

// 	// if err2 != nil {
// 	// 	panic(err2.Error())
// 	// }
// 	// defer insert2.Close()

// 	results, err := db.Query("SELECT * FROM courselist")

// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	for results.Next() {
// 		var n item
// 		err := results.Scan(&n.name, &n.department, &n.crn)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		fmt.Println(n)
// 	}
// 	defer results.Close()
// }
