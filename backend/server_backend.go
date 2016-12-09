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

type PhotoList struct {
    Photos       []*Photo `json="photos"`
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
    router.GET("/photo/latest/:page", l.getLatestPhotos)

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

    rows, res,  err := db.Query("select * from hat4cat.photos where photodI=%d", id)
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
    var photos []*Photo
    db := l.connectToDB()
    page, err := strconv.Atoi(ps.ByName("page"))
    checkError(w, err)
    l1 := page * 11 - 11
    if l1 < 0 { l1 = 0}
    l2 := l1 + 11

    rows, res,  err := db.Query("select * from hat4cat.photos order by photoId desc limit %d, %d", l1, l2)
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
            photo := &Photo{row.Int(id), row.Str(imgName), row.Str(imgDesc), row.Str(image), row.Str(created), row.Int(uid)}
            photos = append(photos, photo)

        }
        photoList := &PhotoList{Photos: photos}

        jsonBody, err := json.Marshal(photoList)
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

