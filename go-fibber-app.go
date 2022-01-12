package main

import (
	"fmt"

	userModel "github.com/jatan/go-fiber-app/model"
	_ "gorm.io/driver/sqlite"
)

func main() {
	fmt.Println("Go ORM Tutorial")

	userModel.InitialMigration()
	// Handle Subsequent requests
	userModel.HandleRequests()
}
