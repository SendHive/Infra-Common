package main

import (
	"fmt"
	"infra-comman/db"
	"log"
)

func main() {
	dbConn, err := db.NewDbRequest()
	if err != nil {
		log.Fatal("the error while creating the database connection: ", err.Error())
	}
	dbConn.InitDB()
	fmt.Println("infra comman service entered")
}
