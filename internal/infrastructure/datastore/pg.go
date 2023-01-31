package datastore

import (
	"api/config"
	"api/internal/domain"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func NewPostgreSQL() *gorm.DB {

	if DB != nil {
		fmt.Println("already")
		return DB
	}
	connectString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s %s",
		config.C.Postgres.HOST,
		config.C.Postgres.PORT,
		config.C.Postgres.USER,
		config.C.Postgres.DBNAME,
		config.C.Postgres.PASS,
		config.C.Postgres.OPTION,
	)
	DB, err := gorm.Open(postgres.Open(connectString), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	DB.AutoMigrate(&domain.User{})


	
	
	


	return DB
}
