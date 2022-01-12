package userModel

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	crypto "github.com/jatan/go-fiber-app/crypto"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string
}

func AllUsers(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var users []User
	db.Find(&users)
	respondWithJSON(w, 200, users)
}

func NewUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New User Endpoint Hit")

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]
	email := vars["email"]

	db.Create(&User{Name: name, Email: email})
	respondWithJSON(w, 200, "New User Successfully Created")
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]

	var user User
	db.Where("name = ?", name).Find(&user)
	db.Delete(&user)

	respondWithJSON(w, 200, "Successfully Deleted User")
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]
	email := vars["email"]

	var user User
	db.Where("name = ?", name).Find(&user)

	user.Email = email

	db.Save(&user)
	respondWithJSON(w, 200, "Successfully Updated User")
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	dat, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, 500, "couldn't read request")
		return
	}
	existingUser := User{}
	user := User{}
	err = json.Unmarshal(dat, &user)
	if err != nil {
		respondWithError(w, 500, "couldn't unmarshal parameters")
		return
	}
	temp, err := crypto.Generate("Test")
	fmt.Println(temp)
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.Where("name = ?", user.Name).Find(&existingUser)
	if existingUser.Name != "" {
		respondWithJSON(w, 400, "User "+user.Name+" Already exisits")
	} else {
		db.Create(&user)
		respondWithJSON(w, 200, " New User "+user.Name+" Successfully Created")
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func HandleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/register", RegisterUser).Methods("POST")
	myRouter.HandleFunc("/users", AllUsers).Methods("GET")
	myRouter.HandleFunc("/user/{name}", DeleteUser).Methods("DELETE")
	myRouter.HandleFunc("/user/{name}/{email}", UpdateUser).Methods("PUT")
	myRouter.HandleFunc("/user/{name}/{email}", NewUser).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func InitialMigration() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&User{})
}
