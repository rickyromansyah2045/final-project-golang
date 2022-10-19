package main

import (
	"final-project-golang/database"
)

func main() {

	db, err := database.ConnectDB()

	if err != nil {
		panic(err.Error())
	}

	_ = db
}
