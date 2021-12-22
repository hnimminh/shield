package main

import (
	"fmt"
	"log"
  	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
)


type Article struct {
	Id string `json:"Id"`
	Title string `json:"Title"`
	Desc string `json:"desc"`
	Content string `json:"content"`
}

// let's declare a global Articles array
// that we can then populate in our main function
// to simulate a database
var Articles []Article

func returnSingleArticle(w http.ResponseWriter, r *http.Request){
	fmt.Println("Endpoint Hit: returnSingleArticle")
	vars := mux.Vars(r)
	key := vars["id"]
	for _, article := range Articles {
		if article.Id == key {
			json.NewEncoder(w).Encode(article)
		}
	}
}

func returnAllArticles(w http.ResponseWriter, r *http.Request){
	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(w).Encode(Articles)
}


func homePage(w http.ResponseWriter, r *http.Request){
	//fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
	a := json{status: 200}
	json.NewEncoder(w).Encode(a)
}


func handleRequests() {
	apiRouter := mux.NewRouter().StrictSlash(true)
	apiRouter.HandleFunc("/", homePage)
	apiRouter.HandleFunc("/articles", returnAllArticles)
	apiRouter.HandleFunc("/articles/{id}", returnSingleArticle)
	log.Fatal(http.ListenAndServe(":8080", apiRouter))
}


func main() {
	Articles = []Article{
		Article{Id: "1", Title: "Hello", Desc: "Hello Article Description", Content: "Hello Article Content"},
		Article{Id: "2", Title: "Bonjour", Desc: "Bonjour Article Description", Content: "Bonjour Article Content"},
	}
	fmt.Println("Rest API v2.0 - Mux Routers")
	handleRequests()
}

// https://tutorialedge.net/golang/creating-restful-api-with-golang/#json
