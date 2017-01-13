package main

import (
	"fmt"
	"github.com/rs/cors"
	"github.com/julienschmidt/httprouter" //https://github.com/julienschmidt/httprouter
	"net/http"
	"strconv"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native" // Native engine
	"encoding/json"
	"log"
	"os"
)


type loginDB struct {
	usr  string
	pswd string
}

type User struct {
	Id          int    `json="id"`
	Firstname   string `json="firstname"`
	Lastname    string `json="lastname"`
	GoogleToken string `json="googletoken"`
	AuthToken   string `json="authtoken"`
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

type Comment struct {
	Cid       int    `json="cid"`
	PhotoId   int    `json="photoId"`
	Comment   string `json="comment"`
	Uid       int    `json="uid"`
	Firstname string `json="firstname"`
	Lastname  string `json="lastname"`
	Timestamp string `json="timestamp"`
}

type Comments struct {
	Comments []*Comment `json="comments"`
}

type Rating struct {
	PhotoId int `json="photoId"`
	Rate    int `json="rate"`
	Uid     int `json="uid"`
}

type RateSum struct {
	PhotoId int `json="photoId"`
	RateSum int `json="rateSum"`
}

type Favorite struct {
	PhotoId int `json="photoId"`
	Uid     int `json="uid"`
}

type Thumbnail struct {
	Id        int    `json="id"`
	ImgName   string `json="imgName"`
	Thumbnail string `json="thumbnail"`
}

type ThumbnailList struct {
	Thumbnails []*Thumbnail `json="thumbnails"`
}

type Toplist struct {
	Toplist []*Thumbnail `json="toplist"`
}

/*****************************************
*** Starts the http-server with the    ***
*** different commands (get, post etc) ***
*****************************************/

func (l *loginDB) startBackend() {
	ip := "130.240.170.62:1026"
	//ip := "localhost:1026"
	router := httprouter.New()
	router.GET("/user/:id", l.getUser)
	router.POST("/user", l.newUser)
	router.POST("/photo", l.savePhoto)
	router.GET("/photo/get/:pid", l.getPhoto)
    router.DELETE("/photo/:pid/:uid", l.deletePhoto)
	router.GET("/photo/user/:uid", l.getUserPhotos)
	router.GET("/photo/favorite/:uid", l.getUserFavoritePhotos)
	router.GET("/photo/latest/:page", l.getLatestPhotos)
	router.GET("/photo/latest/", l.getLatestPhotos)
	router.GET("/photo/top/:list", l.getToplist)
	router.GET("/comments/:pid", l.getComments)
	router.POST("/comment", l.newComment)
	router.POST("/rating", l.newRating)
	router.PUT("/rating/update", l.updateRating)
	router.GET("/rating/:pid/:uid", l.getRating)
	router.GET("/rating/:pid", l.getRatingSum)
	router.POST("/favorite", l.addFavorite)
	router.DELETE("/favorite/:pid/:uid", l.removeFavorite)
	router.GET("/favorite/:pid/:uid", l.getFavorite)


	//handler := cors.Default().Handler(router)

	c := cors.New(cors.Options{
    	AllowedMethods: []string{"POST", "GET", "DELETE", "PUT"},
	})

	// Insert the middleware
	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(ip, handler))
	fmt.Println("running on", ip)

	//log.Fatal(http.ListenAndServe("localhost:1026", handler))
	//fmt.Println("running on localhost:1026")
}


/*******************************************************
*************** USER HANDLERS *************************
*******************************************************/

func (l *loginDB) newUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("POST newUser")
	db := l.connectToDB()
	dec := json.NewDecoder(r.Body)
	user := User{}
	err := dec.Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	token, err := strconv.Atoi(user.GoogleToken)
	if statusCheck(user.AuthToken) {
		_, res, err := db.Query("select count(*) as noUsers from users where googletoken=%d", token)
		checkError(w, err)
		n := res.Map("noUsers")
		fmt.Println(n)
		if n == 0 {
			res, err := db.Prepare("insert into users (firstname, lastname, googletoken) values (?, ?, ?)")
			checkError(w, err)
			_, err = res.Run(user.Firstname, user.Lastname, user.GoogleToken)
			checkError(w, err)
		}

	} else {
		//	fmt.Println("token is not valid")
	}

}

func (l *loginDB) getUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("GET getUser")
    var rows []mysql.Row
    var res mysql.Result
    var err error
	db := l.connectToDB()
    input := ps.ByName("id")

    if len(input) < 12 {
        id, err := strconv.Atoi(ps.ByName("id"))
        checkError(w, err)
        rows, res, err = db.Query("select * from users where uid=%d", id)
        checkError(w, err)
    } else {
        rows, res, err = db.Query("select * from users where googletoken=%s", input)
        checkError(w, err)
    }


	if rows == nil {
		w.WriteHeader(404)
	} else {
		for _, row := range rows {
			id := res.Map("uid")
			firstname := res.Map("firstname")
			lastname := res.Map("lastname")
			googletoken := res.Map("googletoken")
			authtoken := ""
			usr := &User{row.Int(id), row.Str(firstname), row.Str(lastname), row.Str(googletoken), authtoken}

			jsonBody, err := json.Marshal(usr)
			w.WriteHeader(200) // is ok
			w.Write(jsonBody)
			checkError(w, err)
		}
	}

}

/*******************************************************
*************** PHOTO HANDLERS *************************
*******************************************************/
func (l *loginDB) savePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	fmt.Println("POST savePhoto")
	db := l.connectToDB()
	dec := json.NewDecoder(r.Body)
	photo := Photo{}
	err := dec.Decode(&photo)
	if err != nil {
		log.Fatal(err)
	}
	res, err := db.Prepare("insert into photos (name, description, image, uid, thumbnail) values (?, ?, ?, ?, ?)")
	checkError(w, err)

	_, err = res.Run(photo.ImgName, photo.ImgDesc, photo.Image, photo.Uid, photo.Thumbnail)
	checkError(w, err)
}

func (l *loginDB) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Println("DELETE deletePhoto")
    db := l.connectToDB()
    pid, err := strconv.Atoi(ps.ByName("pid"))
    checkError(w, err)

    uid, err := strconv.Atoi(ps.ByName("uid"))
    checkError(w, err)
    
    _, _, err = db.Query("delete from photos where photoId=%d and uid=%d", pid, uid)
    checkError(w, err)
}

func (l *loginDB) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("GET getPhoto")
	db := l.connectToDB()
	id, err := strconv.Atoi(ps.ByName("pid"))
	checkError(w, err)

	rows, res, err := db.Query("select * from photos where photoId=%d", id)
	checkError(w, err)

	if rows == nil {
		w.WriteHeader(404)
	} else {
		for _, row := range rows {
			id := res.Map("photoId")
			imgName := res.Map("name")
			imgDesc := res.Map("description")
			image := res.Map("image")
			created := res.Map("date")
			uid := res.Map("uid")
			thumbnail := res.Map("thumbnail")
			photo := &Photo{row.Int(id), row.Str(imgName), row.Str(imgDesc), row.Str(image), row.Str(created), row.Int(uid), row.Str(thumbnail)}

			jsonBody, err := json.Marshal(photo)
			w.Write(jsonBody)
			checkError(w, err)
		}
	}
}

func (l *loginDB) getLatestPhotos(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("GET getLatestPhotos")
	var thumbnails []*Thumbnail
	p := ps.ByName("page")
	if len(p) == 0 {
		p = "1"
	}
	limit := 47
	db := l.connectToDB()
	page, err := strconv.Atoi(p)
	checkError(w, err)


	l1 := page*limit - limit
	if l1 < 0 {
		l1 = 0
	}
	l2 := l1 + limit

	rows, res, err := db.Query("select photoId, name, thumbnail from photos order by photoId desc limit %d, %d", l1, l2)
	checkError(w, err)

	if rows == nil {
		w.WriteHeader(404)
	} else {
		for _, row := range rows {
			id := res.Map("photoId")
			imgName := res.Map("name")
			thumbnail := res.Map("thumbnail")
			tn := &Thumbnail{row.Int(id), row.Str(imgName), row.Str(thumbnail)}
			thumbnails = append(thumbnails, tn)

		}
		tnList := &ThumbnailList{Thumbnails: thumbnails}

		jsonBody, err := json.Marshal(tnList)
		w.Write(jsonBody)
		checkError(w, err)
	}
}

func (l *loginDB) getUserPhotos(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("GET getUserPhotos")
	var thumbnails []*Thumbnail
	db := l.connectToDB()
	uid, err := strconv.Atoi(ps.ByName("uid"))
	checkError(w, err)

	rows, res, err := db.Query("select photoId, name, thumbnail from photos where uid=%d order by photoId desc", uid)
	checkError(w, err)

	if rows == nil {
		w.WriteHeader(404)
	} else {
		for _, row := range rows {
			id := res.Map("photoId")
			imgName := res.Map("name")
			thumbnail := res.Map("thumbnail")
			tn := &Thumbnail{row.Int(id), row.Str(imgName), row.Str(thumbnail)}
			thumbnails = append(thumbnails, tn)

		}
		tnList := &ThumbnailList{Thumbnails: thumbnails}

		jsonBody, err := json.Marshal(tnList)
		w.Write(jsonBody)
		checkError(w, err)
	}
}

func (l *loginDB) getUserFavoritePhotos(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("GET getUserFavoritePhotos")
	var thumbnails []*Thumbnail
	db := l.connectToDB()
	uid, err := strconv.Atoi(ps.ByName("uid"))
	checkError(w, err)

	rows, res, err := db.Query("select photoId, name, thumbnail from photos as p right join favorite as f on f.pid = p.photoId where f.uid=%d order by photoId desc", uid)
	checkError(w, err)

	if rows == nil {
		w.WriteHeader(404)
	} else {
		for _, row := range rows {
			id := res.Map("photoId")
			imgName := res.Map("name")
			thumbnail := res.Map("thumbnail")
			tn := &Thumbnail{row.Int(id), row.Str(imgName), row.Str(thumbnail)}
			thumbnails = append(thumbnails, tn)

		}
		tnList := &ThumbnailList{Thumbnails: thumbnails}

		jsonBody, err := json.Marshal(tnList)
		w.Write(jsonBody)
		checkError(w, err)
	}
}

func (l *loginDB) getToplist(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("GET getToplist")
	var rows []mysql.Row
	var res mysql.Result
	var err error
	list := ps.ByName("list")
	var top []*Thumbnail
	db := l.connectToDB()

	if list == "rate" {
		rows, res, err = db.Query("select p.photoId, p.name, p.thumbnail, sum(case when rate is null then 0 else rate end) as ratingSum from rating as r right join hat4cat.photos as p on p.photoId=r.photoId group by photoId order by ratingsum desc limit 0,9")
		checkError(w, err)
	} else if list == "comment" {
		rows, res, err = db.Query("select p.photoId, p.name, p.thumbnail, count(comment) as noComments from comment as c right join photos as p on p.photoId=c.photoId group by photoId order by noComments desc limit 0,9")
		checkError(w, err)
	} else if list == "favorite" {
		rows, res, err = db.Query("select p.photoId, p.name, p.thumbnail, count(*) as noFavorites from favorite as f left join photos as p on p.photoId=f.pid group by photoId order by noFavorites desc limit 0,9")
		checkError(w, err)
	} else {
		fmt.Println("Forbidden getToplist request")
		w.WriteHeader(404)
		return
	}

	if rows == nil {
		w.WriteHeader(404)
	} else {
		for _, row := range rows {
			id := res.Map("photoId")
			imgName := res.Map("name")
			thumbnail := res.Map("thumbnail")
			tn := &Thumbnail{row.Int(id), row.Str(imgName), row.Str(thumbnail)}
			top = append(top, tn)

		}
		topList := &Toplist{Toplist: top}

		jsonBody, err := json.Marshal(topList)
		w.Write(jsonBody)
		checkError(w, err)
	}
}

/*******************************************************
*************** COMMENT HANDLERS ***********************
*******************************************************/

func (l *loginDB) newComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("POST newComment")
	db := l.connectToDB()
	dec := json.NewDecoder(r.Body)
	comment := Comment{}
	err := dec.Decode(&comment)
	if err != nil {
		log.Fatal(err)
	}
	res, err := db.Prepare("insert into comment (photoId, comment, uid) values (?, ?, ?)")
	checkError(w, err)

	_, err = res.Run(comment.PhotoId, comment.Comment, comment.Uid)
	checkError(w, err)
}

func (l *loginDB) getComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("GET getComments")
	var comments []*Comment
	db := l.connectToDB()

	id, err := strconv.Atoi(ps.ByName("pid"))
	checkError(w, err)

	rows, res, err := db.Query("select cid, photoId, comment, c.uid, firstname, lastname, timestamp  from comment as c join users as u on c.uid = u.uid where photoId=%d order by cid desc", id)
	checkError(w, err)

	if rows == nil {
		fmt.Println("No comments found for photoId:", id)
	} else {
		for _, row := range rows {
			cid := res.Map("cid")
			photoId := res.Map("photoId")
			comment := res.Map("comment")
			uid := res.Map("uid")
			fname := res.Map("firstname")
			lname := res.Map("lastname")
			timestamp := res.Map("timestamp")
			c := &Comment{row.Int(cid), row.Int(photoId), row.Str(comment), row.Int(uid), row.Str(fname), row.Str(lname), row.Str(timestamp)}
			comments = append(comments, c)

		}
		commentL := &Comments{Comments: comments}

		jsonBody, err := json.Marshal(commentL)
		w.WriteHeader(200) // is ok
		w.Write(jsonBody)
		checkError(w, err)
	}
}

/*******************************************************
**************** RATING HANDLERS ***********************
*******************************************************/

func (l *loginDB) newRating(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("POST newRating")
	db := l.connectToDB()
	dec := json.NewDecoder(r.Body)
	rating := Rating{}
	err := dec.Decode(&rating)
	if err != nil {
		log.Fatal(err)
	}
	res, err := db.Prepare("insert into rating (photoId, rate, uid) values (?, ?, ?)")
	checkError(w, err)

	_, err = res.Run(rating.PhotoId, rating.Rate, rating.Uid)
	checkError(w, err)
}

func (l *loginDB) updateRating(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("PUT updateRating")
	db := l.connectToDB()
	dec := json.NewDecoder(r.Body)
	rating := Rating{}
	err := dec.Decode(&rating)
	if err != nil {
		log.Fatal(err)
	}
	res, err := db.Prepare("update rating set rate=? where photoId=? and uid=?")
	checkError(w, err)

	_, err = res.Run(rating.Rate, rating.PhotoId, rating.Uid)
	checkError(w, err)
}

func (l *loginDB) getRating(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("GET getRating")
	db := l.connectToDB()

	pid, err := strconv.Atoi(ps.ByName("pid"))
	checkError(w, err)

	uid, err := strconv.Atoi(ps.ByName("uid"))
	checkError(w, err)

	rows, res, err := db.Query("select photoId, rate, uid from rating where photoId=%d and uid=%d", pid, uid)
	checkError(w, err)

	if rows == nil {
		fmt.Println("No ratings found for photoId:", pid)
	} else {
		for _, row := range rows {
			photoId := res.Map("photoId")
			rate := res.Map("rate")
			uid := res.Map("uid")
			rating := &Rating{row.Int(photoId), row.Int(rate), row.Int(uid)}
			jsonBody, err := json.Marshal(rating)
			w.WriteHeader(200) // is ok
			w.Write(jsonBody)
			checkError(w, err)
		}

	}
}

func (l *loginDB) getRatingSum(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("GET getRatingSum")
	db := l.connectToDB()

	pid, err := strconv.Atoi(ps.ByName("pid"))
	checkError(w, err)

	rows, res, err := db.Query("select photoId, sum(rate) as ratingSum from rating where photoId=%d", pid)
	checkError(w, err)

	if rows == nil {
		w.WriteHeader(404)
	} else {
		for _, row := range rows {
			photoId := res.Map("photoId")
			rateS := res.Map("ratingSum")
			rateSum := &RateSum{row.Int(photoId), row.Int(rateS)}

			jsonBody, err := json.Marshal(rateSum)
			w.WriteHeader(200) // is ok
			w.Write(jsonBody)
			checkError(w, err)
		}
	}
}

/*******************************************************
**************** FAVORITE HANDLERS ***********************
*******************************************************/

func (l *loginDB) addFavorite(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("POST addFavorite")
	db := l.connectToDB()
	dec := json.NewDecoder(r.Body)
	fav := Favorite{}
	err := dec.Decode(&fav)
	if err != nil {
		log.Fatal(err)
	}
	res, err := db.Prepare("insert into favorite (pid, uid) values (?, ?)")
	checkError(w, err)

	_, err = res.Run(fav.PhotoId, fav.Uid)
	checkError(w, err)
}


func (l *loginDB) removeFavorite(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Println("DELETE removeFavorite")
    db := l.connectToDB()
    pid, err := strconv.Atoi(ps.ByName("pid"))
    checkError(w, err)

    uid, err := strconv.Atoi(ps.ByName("uid"))
    checkError(w, err)
    
    _, _, err = db.Query("delete from favorite where pid=%d and uid=%d", pid, uid)
    checkError(w, err)
}

func (l *loginDB) getFavorite(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("GET getFavorite")
	db := l.connectToDB()

	pid, err := strconv.Atoi(ps.ByName("pid"))
	checkError(w, err)

	uid, err := strconv.Atoi(ps.ByName("uid"))
	checkError(w, err)

	rows, res, err := db.Query("select pid, uid from favorite where pid=%d and uid=%d", pid, uid)
	checkError(w, err)

	if rows == nil {
	} else {
		for _, row := range rows {
			photoId := res.Map("pid")
			uid := res.Map("uid")
			fav := &Favorite{row.Int(photoId), row.Int(uid)}
			jsonBody, err := json.Marshal(fav)
			w.WriteHeader(200) // is ok
			w.Write(jsonBody)
			checkError(w, err)
		}

	}
}

/*******************************************************
*************** ERROR HANDLERS *************************
*******************************************************/

func checkError(w http.ResponseWriter, err error) {
	if err != nil {
		w.WriteHeader(500) // error
		fmt.Println(err)
		fmt.Fprintf(w, "Bad input")
		w.Write([]byte("{ error: bad input }"))
	}
}

/*******************************************************
*************** DATABASE HANDLERS **********************
*******************************************************/

func (l *loginDB) connectToDB() mysql.Conn {
	ip_db := "130.240.170.62:3306"
	db_name := "hat4cat"
	db := mysql.New("tcp", "", ip_db, l.usr, l.pswd, db_name)
	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to database")
	return db
}

//check the status code of the loged in user
func statusCheck(token string) bool {
	auth_token := "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + token
	resp, err := http.Get(auth_token)
	if (err == nil) && (resp.StatusCode == 200) {
		return true
	} else {
		return false
	}
}

/*******************************************************
********************* MAIN *****************************
*******************************************************/

func main() {
	usr := os.Args[1]
	pswd := os.Args[2]
	login := &loginDB{usr, pswd}
	login.startBackend()
}
