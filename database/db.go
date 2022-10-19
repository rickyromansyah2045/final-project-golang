package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	HOST_POSTGRES = "localhost"
	PORT_POSTGRES = 5432
	DB_POSTGRES   = "final_project_go"
	USER_POSTGRES = "postgres"
	PASS_POSTGRES = "Noerhick_02"
)

var (
	db  *gorm.DB
	err error
)

func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST_POSTGRES, PORT_POSTGRES, USER_POSTGRES, PASS_POSTGRES, DB_POSTGRES,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate()

	return db
}
