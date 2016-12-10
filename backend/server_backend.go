package main

import (
	"fmt"
    "github.com/rs/cors"
    //"encoding/json"
    //"time"
    "html/template"
    "strconv"
	"net/http"
	"github.com/julienschmidt/httprouter" //https://github.com/julienschmidt/httprouter
	//"github.com/ziutek/mymysql"
    "github.com/ziutek/mymysql/mysql"
    _ "github.com/ziutek/mymysql/native" // Native engine
    // _ "github.com/ziutek/mymysql/thrsafe" // Thread safe engine
    //"database/sql"
    //"bufio"
	"os"
	"log"
    "encoding/json"
    //b64 "encoding/base64"
)

type loginDB struct {
    usr     string
    pswd     string
}

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


type Comment struct {
    Cid         int `json="cid"`
    PhotoId     int `json="photoId"`
    Comment     string `json="comment"`
    Uid         int `json="uid"`
    Timestamp   string `json="timestamp"`
}

type Comments struct {
    Comments       []*Comment `json="comments"`
}

type Rating struct {
    Rid         int `json="rid"`
    PhotoId     int `json="photoId"`
    Rate        int `json="rate"`
}

type Thumbnail struct {
    Id          int `json="id"`
    ImgName     string `json="imgName"`
    Thumbnail   string `json="thumbnail"`
}

type ThumbnailList struct {
    Thumbnails       []*Thumbnail `json="thumbnails"`
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
    router.GET("/photo/:id", l.getPhoto)
    router.GET("/latest/:page", l.getLatestPhotos)
    router.GET("/comments/:id", l.getComments)
    router.POST("/comment", l.newComment)

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
    db := l.connectToDB()
    dec := json.NewDecoder(r.Body)
    user := User{}
    err := dec.Decode(&user)
    if err != nil {
        log.Fatal(err)
    }
    res,  err := db.Prepare("insert into hat4cat.users (firstname, lastname) values (?, ?)")
    checkError(w, err)

    _, err = res.Run(user.Firstname, user.Lastname)
    checkError(w, err)
}


func (l *loginDB) getUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    db := l.connectToDB()

    id, err := strconv.Atoi(ps.ByName("id"))
    checkError(w, err)

    rows, res,  err := db.Query("select * from hat4cat.users where uid=%d", id)
    checkError(w, err)

    if rows == nil {
        w.WriteHeader(404)
    } else {
        for _, row := range rows {
        id := res.Map("uid")
        firstname := res.Map("firstname")
        lastname := res.Map("lastname")
        usr := &User{row.Int(id), row.Str(firstname), row.Str(lastname)}

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
    db := l.connectToDB()
    dec := json.NewDecoder(r.Body)
    photo := Photo{}
    err := dec.Decode(&photo)
    if err != nil {
        log.Fatal(err)
    }
    res,  err := db.Prepare("insert into hat4cat.photos (name, description, image, uid, thumbnail) values (?, ?, ?, ?, ?)")
    checkError(w, err)

    _, err = res.Run(photo.ImgName, photo.ImgDesc, photo.Image, photo.Uid, photo.Thumbnail)
    checkError(w, err)
}


func (l *loginDB) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    db := l.connectToDB()
    id, err := strconv.Atoi(ps.ByName("id"))
    checkError(w, err)

    rows, res,  err := db.Query("select * from hat4cat.photos where photoId=%d", id)
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
    var thumbnails []*Thumbnail
    limit := 23
    db := l.connectToDB()
    page, err := strconv.Atoi(ps.ByName("page"))
    checkError(w, err)
    l1 := page * limit - limit
    if l1 < 0 { l1 = 0}
    l2 := l1 + limit

    rows, res,  err := db.Query("select photoId, name, thumbnail from hat4cat.photos order by photoId desc limit %d, %d", l1, l2)
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

/*******************************************************
*************** COMMENT HANDLERS ***********************
*******************************************************/

func (l *loginDB) newComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    db := l.connectToDB()
    dec := json.NewDecoder(r.Body)
    comment := Comment{}
    err := dec.Decode(&comment)
    if err != nil {
        log.Fatal(err)
    }
    res,  err := db.Prepare("insert into hat4cat.comment (photoId, comment, uid) values (?, ?, ?)")
    checkError(w, err)

    _, err = res.Run(comment.PhotoId, comment.Comment, comment.Uid)
    checkError(w, err)
}


func (l *loginDB) getComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    var comments []*Comment
    db := l.connectToDB()

    id, err := strconv.Atoi(ps.ByName("id"))
    checkError(w, err)

    rows, res,  err := db.Query("select * from hat4cat.comment where photoId=%d", id)
    checkError(w, err)

    if rows == nil {
        w.WriteHeader(404)
    } else {
        for _, row := range rows {
            cid := res.Map("cid")
            photoId := res.Map("photoId")
            comment := res.Map("comment")
            uid := res.Map("uid")
            timestamp := res.Map("timestamp")
            c := &Comment{row.Int(cid), row.Int(photoId), row.Str(comment), row.Int(uid), row.Str(timestamp)}
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

/*******************************************************
********************* MAIN *****************************
*******************************************************/

func main() {
    usr := os.Args[1]
    pswd := os.Args[2]
    login := &loginDB{usr, pswd}
    login.startBackend()
}

