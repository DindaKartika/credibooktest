package main

import (
	"fmt"

	"credibooktest/config"
	"credibooktest/models"
	"credibooktest/router"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	fmt.Println("Welcome to credibook apps")

	e := router.New()

	// init database
	autoCreateTables()
	// autoMigrateTables()

	e.Start(":8000")
}

// autoCreateTables: create database tables using GORM
// will be moved to database/seeder
func autoCreateTables() {
	if !config.App.DBConfig.HasTable(&models.User{}) {
		config.App.DBConfig.CreateTable(&models.User{})
		if config.App.ENV == "dev" {
			var user = models.User{
				Username: "admin",
				Password: "admin",
				IsAdmin:  true,
			}
			password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
			user.Password = string(password)
			config.App.DBConfig.Create(&user)
		}
	}
	if !config.App.DBConfig.HasTable(&models.Transaction{}) {
		config.App.DBConfig.CreateTable(&models.Transaction{})
		if config.App.ENV == "dev" {
			var transaction = models.Transaction{
				UserID: 1,
				Amount: 100000,
				Notes:  "Belanja",
				Type:   "expense",
			}
			config.App.DBConfig.Create(&transaction)
		}
	}
}

// autoMigrateTables: migrate table columns using GORM
// will be moved to database/migration
// func autoMigrateTables() {
// 	config.App.DBConfig.AutoMigrate(&models.User{})
// 	config.App.DBConfig.AutoMigrate(&models.Transaction{})
// }
