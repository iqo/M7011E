package main

import (
	"fmt"
	"html/template"
	"net/http"
    "log"
    "os"
    "strings"
    "encoding/json"
    "strconv"
    //"io/ioutil"
	//"github.com/julienschmidt/httprouter" //https://github.com/julienschmidt/httprouter

)
type User struct {
    Id          int `json="id"`
    Firstname   string `json="firstname"`
    Lastname    string `json="lastname"`
    }

type Photo struct {
    Id          int `json="id"`
    ImgName     string `json="imgName"`
    ImgDesc     string `json="imgDesc"`
    Image       string `json="image"`
    Created     string `json="created"`
    Uid         int `json="uid"`
    Thumbnail   string `json="thumbnail"`
}

type PhotoView struct {
    Id          int `json="id"`
    ImgName     string `json="imgName"`
    ImgDesc     string `json="imgDesc"`
    Image       string `json="image"`
    Created     string `json="created"`
    Uid         int `json="uid"`
    Firstname   string `json="firstname"`
    Lastname    string `json="lastname"`
    RatingSum   int `json="ratingSum"`
}

type RateSum struct {
    PhotoId     int `json="photoId"`
    RateSum     int `json="rateSum"`
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
    pResponse, err := http.Get("http://130.240.170.62:1026/photo/" + id[2])
    checkError(w, err)
    defer pResponse.Body.Close()
    dec := json.NewDecoder(pResponse.Body)
    photo := Photo{}
    err = dec.Decode(&photo)

    checkError(w, err)
    uResponse, err := http.Get("http://130.240.170.62:1026/user/" + strconv.Itoa(photo.Uid))
    checkError(w, err)
    defer uResponse.Body.Close()
    dec = json.NewDecoder(uResponse.Body)
    user := User{}
    err = dec.Decode(&user)

    checkError(w, err)
    rResponse, err := http.Get("http://130.240.170.62:1026/rating/" + id[2])
    checkError(w, err)
    defer rResponse.Body.Close()
    dec = json.NewDecoder(rResponse.Body)
    r := RatingSum{}
    err = dec.Decode(&r)

    checkError(w, err)
    photoV := &PhotoView{photo.Id, photo.ImgName, photo.ImgDesc, photo.Image, photo.Created, photo.Uid, user.Firstname, user.Lastname, r.RateSum}
    t :=template.Must(template.ParseFiles("static/index.html", "static/templates/photo.tmp"))
    t.Execute(w, photoV)
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

