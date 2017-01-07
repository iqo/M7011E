package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"io"
	//"github.com/julienschmidt/httprouter" //https://github.com/julienschmidt/httprouter
)

type User struct {
	Id        int    `json="id"`
	Firstname string `json="firstname"`
	Lastname  string `json="lastname"`
}

type UserView struct {
    Id          int    `json="id"`
    Firstname   string `json="firstname"`
    Lastname    string `json="lastname"`
    Thumbnails  *ThumbnailList `json="thumbnails"`
    Favorites   *ThumbnailList `json="favorites"`
}

type Photo struct {
	Id        int    `json="id"`
	ImgName   string `json="imgName"`
	ImgDesc   string `json="imgDesc"`
	Image     string `json="image"`
	Created   string `json="created"`
	Uid       int    `json="uid"`
	Thumbnail string `json="thumbnail"`
}

type PhotoView struct {
	Id        int    `json="id"`
	ImgName   string `json="imgName"`
	ImgDesc   string `json="imgDesc"`
	Image     string `json="image"`
	Created   string `json="created"`
	Uid       int    `json="uid"`
	Firstname string `json="firstname"`
	Lastname  string `json="lastname"`
	RatingSum int    `json="ratingSum"`
}

type ToplistView struct {
    Heading     string `json="heading"`
    Text        string `json="text"`
    DivId       string `json="divId"`
    Toplist     *Toplist `json="toplist"`

}

type TextView struct {
    Heading     string `json="heading"`
    Text        string `json="text"`
}

type ThumbnailList struct {
    Thumbnails []*Thumbnail `json="thumbnails"`
}

type Toplist struct {
    Toplist     []*Thumbnail `json="toplist"`
}

type Thumbnail struct {
    Id        int    `json="id"`
    ImgName   string `json="imgName"`
    Thumbnail string `json="thumbnail"`
}

type RateSum struct {
	PhotoId int `json="photoId"`
	RateSum int `json="rateSum"`
}

/*****************************************
*** Adds content on website            ***
*****************************************/
func IndexHandler(w http.ResponseWriter, r *http.Request) {
    heading := "Welcome to Hats for cats!"
    text := "On this page you can put funny hats on your cat pictures! Create an account and upload your favorite pictures of your (or someone else's) cat today. Get some inspiration from the latest cats below."
	tv := &TextView{heading, text}
    template.Must(template.ParseFiles("static/index.html", "static/templates/latestCats.tmp")).Execute(w, tv)
}

func LatestPhotosHandler(w http.ResponseWriter, r *http.Request) {
    heading := "Latest cats in hats"
    text := "The most recent cats in hats are placed here. Click on them and let everyone know what you think of their masterpiece."
    tv := &TextView{heading, text}
    template.Must(template.ParseFiles("static/index.html", "static/templates/latestCats.tmp")).Execute(w, tv)

}



func ToplistRateHandler(w http.ResponseWriter, r *http.Request) {
    var topL []*Thumbnail
    response, err := http.Get("http://130.240.170.62:1026/photo/top/rate")
    checkError(w, err)
    defer response.Body.Close()
    dec := json.NewDecoder(response.Body)
    toplist := Toplist{}
    //thumbnail := Thumbnail{}
    err = dec.Decode(&toplist)
    checkError(w, err)
    for _, t := range toplist.Toplist {
        tn := &Thumbnail{t.Id, t.ImgName, t.Thumbnail}
        topL = append(topL, tn)
    }
    tl := &Toplist{Toplist: topL}

    heading := "Alltime mjaow"
    text := "The cats with the highest rating score of all time is placed here. Why top 9 instead of 10 you ask? That's because cats have 9 lives silly and this makes sence!"
    divId := "toplist-alltime"

    toplistV := &ToplistView{Heading: heading, Text: text, DivId: divId, Toplist: tl}

    t := template.Must(template.ParseFiles("static/index.html", "static/templates/toplist.tmp"))
    t.Execute(w, toplistV)
}

func ToplistCommentHandler(w http.ResponseWriter, r *http.Request) {
    var topL []*Thumbnail
    response, err := http.Get("http://130.240.170.62:1026/photo/top/comment")
    checkError(w, err)
    defer response.Body.Close()
    dec := json.NewDecoder(response.Body)
    toplist := Toplist{}
    //thumbnail := Thumbnail{}
    err = dec.Decode(&toplist)
    checkError(w, err)
    for _, t := range toplist.Toplist {
        tn := &Thumbnail{t.Id, t.ImgName, t.Thumbnail}
        topL = append(topL, tn)
    }
    tl := &Toplist{Toplist: topL}

    heading := "Highest mjaow"
    text := "The cats with the most comments of all time is placed here. Why top 9 instead of 10 you ask? That's because cats have 9 lives silly and this makes sence!"
    divId := "toplist-comment"

    toplistV := &ToplistView{Heading: heading, Text: text, DivId: divId, Toplist: tl}

    t := template.Must(template.ParseFiles("static/index.html", "static/templates/toplist.tmp"))
    t.Execute(w, toplistV)
}

func ToplistFavoriteHandler(w http.ResponseWriter, r *http.Request) {
    var topL []*Thumbnail
    var tl *Toplist
    response, err := http.Get("http://130.240.170.62:1026/photo/top/favorite")
    checkError(w, err)
    defer response.Body.Close()
    dec := json.NewDecoder(response.Body)
    toplist := Toplist{}
    //thumbnail := Thumbnail{}
    err = dec.Decode(&toplist)
    if err != io.EOF && len(toplist.Toplist) != 0 {
        checkError(w, err)
        for _, t := range toplist.Toplist {
            tn := &Thumbnail{t.Id, t.ImgName, t.Thumbnail}
            topL = append(topL, tn)
        }
        tl = &Toplist{Toplist: topL}
    } else {
        tl = &Toplist{Toplist: nil}
    }

    heading := "Most lovable cats"
    text := "These are the cats that most people have added as one of their favorites. Why top 9 instead of 10 you ask? That's because cats have 9 lives silly and this makes sence!"
    divId := "toplist-favorite"

    toplistV := &ToplistView{Heading: heading, Text: text, DivId: divId, Toplist: tl}

    t := template.Must(template.ParseFiles("static/index.html", "static/templates/toplist.tmp"))
    t.Execute(w, toplistV)
}



func PhotoHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.URL.Path, "/")
    if id[2] == "" {
        http.Redirect(w, r, "/", 301)
    } else {
    	pResponse, err := http.Get("http://130.240.170.62:1026/photo/get/" + id[2])
    	checkError(w, err)
    	defer pResponse.Body.Close()
    	dec := json.NewDecoder(pResponse.Body)
    	photo := Photo{}
    	err = dec.Decode(&photo)
        if !checkEOF(err, w, r){

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
        	rate := RateSum{}
        	err = dec.Decode(&rate)

        	checkError(w, err)
        	photoV := &PhotoView{photo.Id, photo.ImgName, photo.ImgDesc, photo.Image, photo.Created, photo.Uid, user.Firstname, user.Lastname, rate.RateSum}
        	t := template.Must(template.ParseFiles("static/index.html", "static/templates/photo.tmp"))
        	t.Execute(w, photoV)
        }
    }
}


func MyPageHandler(w http.ResponseWriter, r *http.Request) {
    uid := strings.Split(r.URL.Path, "/")
    if uid[2] == "" || uid[2] == "0"{
        // if no input
        http.Redirect(w, r, "/", 301)
    } else {
        var pl *ThumbnailList
        var userFav *ThumbnailList

        response, err := http.Get("http://130.240.170.62:1026/photo/user/" + uid[2])
        checkError(w, err)
        defer response.Body.Close()
        dec := json.NewDecoder(response.Body)
        thumbnail := ThumbnailList{}
        err = dec.Decode(&thumbnail)
        pl = thumbnailLoop(&thumbnail, w, err)

        fresponse, err := http.Get("http://130.240.170.62:1026/photo/favorite/" + uid[2])
        checkError(w, err)
        defer fresponse.Body.Close()
        dec = json.NewDecoder(fresponse.Body)
        thumbnail = ThumbnailList{}
        err = dec.Decode(&thumbnail)
        userFav = thumbnailLoop(&thumbnail, w, err)

        uResponse, err := http.Get("http://130.240.170.62:1026/user/" + uid[2])
        checkError(w, err)

        defer uResponse.Body.Close()
        dec = json.NewDecoder(uResponse.Body)
        user := User{}
        err = dec.Decode(&user)

        checkError(w, err)
        userV := &UserView{Id: user.Id, Firstname: user.Firstname, Lastname: user.Lastname, Thumbnails: pl, Favorites: userFav}
        t := template.Must(template.ParseFiles("static/index.html", "static/templates/mypage.tmp"))
        t.Execute(w, userV)
    }
}

func thumbnailLoop(thumbnail *ThumbnailList, w http.ResponseWriter, err error) *ThumbnailList {
    var topL []*Thumbnail
    var pl *ThumbnailList
    if err != io.EOF && len(thumbnail.Thumbnails) != 0 {
            fmt.Println(thumbnail.Thumbnails)
            checkError(w, err)
            for _, t := range thumbnail.Thumbnails {
                tn := &Thumbnail{t.Id, t.ImgName, t.Thumbnail}
                topL = append(topL, tn)
            }
            pl = &ThumbnailList{Thumbnails: topL}
    } else {
            pl = &ThumbnailList{Thumbnails: nil}
        }
        return pl
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseFiles("static/index.html", "static/templates/about.tmp")).Execute(w, nil)
}

func CatMagicHandler(w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseFiles("static/index.html", "static/templates/catmagic.tmp", "static/templates/hats.tmp")).Execute(w, nil)
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
    http.HandleFunc("/toplist/rate", ToplistRateHandler)
    http.HandleFunc("/toplist/comment", ToplistCommentHandler)
    http.HandleFunc("/toplist/favorite", ToplistFavoriteHandler)
    http.HandleFunc("/toplist/latest", LatestPhotosHandler)
	//http.HandleFunc("/toplist", TopListHandler)
	http.HandleFunc("/photo/", PhotoHandler)
	http.HandleFunc("/login", LoginHandler)
    http.HandleFunc("/mypage/", MyPageHandler)

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

func checkEOF(err error, w http.ResponseWriter, r *http.Request) bool {
    if err == io.EOF {
        http.Redirect(w, r, "/", 301)
        return true
    } else {
        return false
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
