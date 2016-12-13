package main

import (
	"fmt"
	"github.com/rs/cors"
	//"encoding/json"
	//"time"
	"github.com/julienschmidt/httprouter" //https://github.com/julienschmidt/httprouter
	"html/template"
	"net/http"
	"strconv"
	//"github.com/ziutek/mymysql"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native" // Native engine
	// _ "github.com/ziutek/mymysql/thrsafe" // Thread safe engine
	//"database/sql"
	//"bufio"
	"encoding/json"
	"log"
	"os"
	//b64 "encoding/base64"
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
	router := httprouter.New()
	router.GET("/", testpage)
	router.GET("/user/:id", l.getUser)
	router.POST("/user", l.newUser)
	router.POST("/photo", l.savePhoto)
	router.GET("/photo/get/:pid", l.getPhoto)
	router.GET("/photo/user/:uid", l.getUserPhotos)
	router.GET("/photo/latest/:page", l.getLatestPhotos)
	router.GET("/photo/top/:list", l.getToplist)
	router.GET("/comments/:id", l.getComments)
	router.POST("/comment", l.newComment)
	router.POST("/rating", l.newRating)
	router.POST("/rating/update", l.updateRating)
	router.GET("/rating/:pid/:uid", l.getRating)
	router.GET("/rating/:pid", l.getRatingSum)

	handler := cors.Default().Handler(router)

	log.Fatal(http.ListenAndServe("130.240.170.62:1026", handler))
	fmt.Println("running on 130.240.170.62:1026")
}

func testpage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	form := "<html><body><form action='/user' method='POST'><p>Firstname: <input type='text' name='firstname' id='fname'/></p><p>Lastname: <input type='text' name='lastname' id='lname' /></p><p><input type='button' onclick = 'testfunc()' value='Testa'/></p></form><script>function testfunc() {var xhr = new XMLHttpRequest(); var fname = document.getElementById('fname').value; var lname = document.getElementById('lname').value;  var user = {}; user.firstname = fname; user.lastname = lname; xhr.open('POST', 'http://localhost:1026/user', true); xhr.setRequestHeader('Content-Type', 'application/json'); xhr.send(JSON.stringify(user))}</script></body></html>"
	t, _ := template.New("webpage").Parse(form)

	t.Execute(w, nil)
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
	if statusCheck(user.AuthToken) {
		rows, _, err := db.Query("select count(*) from hat4cat.users where googletoken=%d", user.GoogleToken)
		checkError(w, err)
		fmt.Println("rows: ", len(rows))
		if len(rows) == 0 {
			res, err := db.Prepare("insert into hat4cat.users (firstname, lastname, googletoken) values (?, ?, ?)")
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
	db := l.connectToDB()

	id, err := strconv.Atoi(ps.ByName("id"))
	checkError(w, err)

	rows, res, err := db.Query("select * from hat4cat.users where uid=%d", id)
	checkError(w, err)

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

func (l *loginDB) getGoogleToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("GET Token")
	db := l.connectToDB()

	id, err := strconv.Atoi(ps.ByName("token"))
	checkError(w, err)

	rows, res, err := db.Query("select * from hat4cat.users where googletoken=%d", id)
	checkError(w, err)

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
	res, err := db.Prepare("insert into hat4cat.photos (name, description, image, uid, thumbnail) values (?, ?, ?, ?, ?)")
	checkError(w, err)

	_, err = res.Run(photo.ImgName, photo.ImgDesc, photo.Image, photo.Uid, photo.Thumbnail)
	checkError(w, err)
}

func (l *loginDB) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("GET getPhoto")
	db := l.connectToDB()
	id, err := strconv.Atoi(ps.ByName("pid"))
	checkError(w, err)

	rows, res, err := db.Query("select * from hat4cat.photos where photoId=%d", id)
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
	limit := 23
	db := l.connectToDB()
	page, err := strconv.Atoi(ps.ByName("page"))
	checkError(w, err)
	l1 := page*limit - limit
	if l1 < 0 {
		l1 = 0
	}
	l2 := l1 + limit

	rows, res, err := db.Query("select photoId, name, thumbnail from hat4cat.photos order by photoId desc limit %d, %d", l1, l2)
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

	rows, res, err := db.Query("select photoId, name, thumbnail from hat4cat.photos where uid=%d order by photoId desc", uid)
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
		rows, res, err = db.Query("select p.photoId, p.name, p.thumbnail, sum(case when rate is null then 0 else rate end) as ratingSum from hat4cat.rating as r right join hat4cat.photos as p on p.photoId=r.photoId group by photoId order by ratingsum desc limit 0,9")
		checkError(w, err)
	} else if list == "comment" {
		rows, res, err = db.Query("select p.photoId, p.name, p.thumbnail, count(comment) as noComments from hat4cat.comment as c right join hat4cat.photos as p on p.photoId=c.photoId group by photoId order by noComments desc limit 0,9")
		checkError(w, err)
	} else {
		fmt.Println("Forbidden getToplist request")
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
	res, err := db.Prepare("insert into hat4cat.comment (photoId, comment, uid) values (?, ?, ?)")
	checkError(w, err)

	_, err = res.Run(comment.PhotoId, comment.Comment, comment.Uid)
	checkError(w, err)
}

func (l *loginDB) getComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("GET getComments")
	var comments []*Comment
	db := l.connectToDB()

	id, err := strconv.Atoi(ps.ByName("id"))
	checkError(w, err)

	rows, res, err := db.Query("select cid, photoId, comment, c.uid, firstname, lastname, timestamp  from hat4cat.comment as c join hat4cat.users as u on c.uid = u.uid where photoId=%d order by cid desc", id)
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
	res, err := db.Prepare("insert into hat4cat.rating (photoId, rate, uid) values (?, ?, ?)")
	checkError(w, err)

	_, err = res.Run(rating.PhotoId, rating.Rate, rating.Uid)
	checkError(w, err)
}

func (l *loginDB) updateRating(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("POST updateRating")
	db := l.connectToDB()
	dec := json.NewDecoder(r.Body)
	rating := Rating{}
	err := dec.Decode(&rating)
	if err != nil {
		log.Fatal(err)
	}
	res, err := db.Prepare("update hat4cat.rating set rate=? where photoId=? and uid=?")
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

	rows, res, err := db.Query("select photoId, rate, uid from hat4cat.rating where photoId=%d and uid=%d", pid, uid)
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

	rows, res, err := db.Query("select photoId, sum(rate) as ratingSum from hat4cat.rating where photoId=%d", pid)
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
*************** ERROR HANDLERS *************************
*******************************************************/

func checkError(w http.ResponseWriter, err error) {
	if err != nil {
		w.WriteHeader(500) // error
		fmt.Println(err)
		fmt.Fprintf(w, "Bad input")
	}
}

/*******************************************************
*************** DATABASE HANDLERS **********************
*******************************************************/

func (l *loginDB) connectToDB() mysql.Conn {
	db := mysql.New("tcp", "", "130.240.170.62:3306", l.usr, l.pswd, "hat4cat")
	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to database")
	return db
}

//ceck the status code of the loged in user
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
