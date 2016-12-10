package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	//"github.com/julienschmidt/httprouter" //https://github.com/julienschmidt/httprouter
)

/*****************************************
*** Adds content on website            ***
*****************************************/
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseFiles("static/index.html", "static/templates/start.html")).Execute(w, nil)
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseFiles("static/index.html", "static/templates/about.html")).Execute(w, nil)
}

func CatMagicHandler(w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseFiles("static/index.html", "static/templates/catmagic.html")).Execute(w, nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseFiles("static/index.html", "static/templates/login.html")).Execute(w, nil)
}

/*****************************************
*** Starts the http-server with the    ***
*** different commands (get, post etc) ***
*****************************************/
func startWebserver(input string) {
	//router := httprouter.New()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	//router.GET("/", IndexHandler)
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/about", AboutHandler)
	http.HandleFunc("/catmagic", CatMagicHandler)
	http.HandleFunc("/login", LoginHandler)

	//var input int
	//fmt.Scan(&input)
	if input == "1" {
		fmt.Println("running on 130.240.170.62:1025")
		log.Fatal(http.ListenAndServe("130.240.170.62:1025", nil))
	} else {
		fmt.Println("running on localhost:1025")
		log.Fatal(http.ListenAndServe("localhost:1025", nil))

	}
}

func main() {
	if len(os.Args) > 1 {
		args := os.Args[1]
		startWebserver(args)
	} else {
		startWebserver("1")
	}
}
