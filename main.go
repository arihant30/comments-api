package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error

type Comment struct {
	ID      uint   `gorm:"primary_key";"AUTO_INCREMENT";json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

func handleRequests() {
	MyRouter := mux.NewRouter().StrictSlash(true)
	MyRouter.HandleFunc("/", HomePage).Methods("GET")
	MyRouter.HandleFunc("/comments", AllComments).Methods("GET")
	MyRouter.HandleFunc("/comments", CreateComment).Methods("POST")
	MyRouter.HandleFunc("/comments/{id}", DeleteComment).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8082", MyRouter))
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Comment API")
}

func CreateComment(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var comment Comment
	json.Unmarshal(reqBody, &comment)
	db.Create(&comment)

	comments := []Comment{}
	db.Find(&comments)
	fmt.Println("Endpoint hit: AllComments")
	json.NewEncoder(w).Encode(comments)
}

func AllComments(w http.ResponseWriter, r *http.Request) {

	comments := []Comment{}
	db.Find(&comments)
	fmt.Println("Endpoint hit: AllComments")
	json.NewEncoder(w).Encode(comments)

}

func DeleteComment(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	var comment Comment

	db.Where("id = ?", id).Delete(&comment)

	fmt.Fprintf(w, "Successfully Deleted Comment \n")

	comments := []Comment{}
	db.Find(&comments)
	fmt.Println("Endpoint hit: AllComments")
	json.NewEncoder(w).Encode(comments)
}

func main() {
	db, err = gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/Comment1?charset=utf8&parseTime=True")
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to db")
	}
	db.AutoMigrate(&Comment{})
	handleRequests()

}
