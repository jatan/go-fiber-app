package main

import (
	"fmt"
	"log"
	"net/http"

	userModel "github.com/jatan/go-fiber-app/model"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "gorm.io/driver/sqlite"
)

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/users", userModel.AllUsers).Methods("GET")
	myRouter.HandleFunc("/user/{name}", userModel.DeleteUser).Methods("DELETE")
	myRouter.HandleFunc("/user/{name}/{email}", userModel.UpdateUser).Methods("PUT")
	myRouter.HandleFunc("/user/{name}/{email}", userModel.NewUser).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func initialMigration() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	// db.AutoMigrate(userModel.(&User){})
}

func main() {
	fmt.Println("Go ORM Tutorial")

	initialMigration()
	// Handle Subsequent requests
	handleRequests()
}
