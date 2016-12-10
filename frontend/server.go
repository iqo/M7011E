package main

import (
	"fmt"
	"html/template"
	"net/http"
    "log"
    "os"
    "strings"
    "io/ioutil"
	//"github.com/julienschmidt/httprouter" //https://github.com/julienschmidt/httprouter

)

type Photo struct {
    Id          int `json="id"`
    ImgName     string `json="imgName"`
    ImgDesc     string `json="imgDesc"`
    Image       string `json="image"`
    Created     string `json="created"`
    Uid         int `json="uid"`
    Thumbnail   string `json="thumbnail"`
}


/*****************************************
*** Adds content on website            ***
*****************************************/
func IndexHandler(w http.ResponseWriter, r *http.Request) {
    template.Must(template.ParseFiles("static/index.html", "static/templates/latestCats.tmp")).Execute(w, nil)
}

func TopListHandler(w http.ResponseWriter, r *http.Request) {
    template.Must(template.ParseFiles("static/index.html", "static/templates/start.html")).Execute(w, nil)
}

func PhotoHandler(w http.ResponseWriter, r *http.Request) {
    id := strings.Split(r.URL.Path, "/")
    response, err := http.Get("http://130.240.170.62:1026/photo/" + id[2])
    fmt.Printf("test")
    checkError(w, err)


    defer response.Body.Close()
    responseData, err := ioutil.ReadAll(response.Body)
    checkError(w, err)
    fmt.Printf(string(responseData))

    //template.Must(template.ParseFiles("static/index.html", "static/templates/start.html")).Execute(w, nil)
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
    template.Must(template.ParseFiles("static/index.html", "static/templates/about.html")).Execute(w, nil)
}

func CatMagicHandler(w http.ResponseWriter, r *http.Request) {
    template.Must(template.ParseFiles("static/index.html", "static/templates/catmagic.html", "static/templates/hats.tmp")).Execute(w, nil)
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
    http.HandleFunc("/toplist", TopListHandler)
    http.HandleFunc("/photo/", PhotoHandler)



    
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

func checkError(w http.ResponseWriter, err error) {
    if err != nil {
        w.WriteHeader(500) // error
        fmt.Println(err)
        fmt.Fprintf(w, "Bad input")
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

