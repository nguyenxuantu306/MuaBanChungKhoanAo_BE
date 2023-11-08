package database

import (
	"mymodule/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
)

var DB *gorm.DB

func Connect() {
	connection, err := gorm.Open(mysql.Open("root:1234@tcp(localhost:3306)/stock?parseTime=true"), &gorm.Config{})

	if err != nil {
		panic("Could not connect to the database")
	}

	DB = connection

	connection.AutoMigrate(&models.Users{})
}
