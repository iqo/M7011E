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
    }


/*****************************************
*** Starts the http-server with the    ***
*** different commands (get, post etc) ***
*****************************************/
func (l *loginDB) startWebserver() {
    router := httprouter.New()
    router.GET("/", testpage)
    router.GET("/user/:id", l.getUser)
    router.POST("/user", l.newUser)
    router.POST("/photo", l.savePhoto)

    handler := cors.Default().Handler(router)

    log.Fatal(http.ListenAndServe("130.240.170.62:1026", handler))
    fmt.Println("running on 130.240.170.62:1026")
}

func testpage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    form := "<html><body><form action='/user' method='POST'><p>Firstname: <input type='text' name='firstname' id='fname'/></p><p>Lastname: <input type='text' name='lastname' id='lname' /></p><p><input type='button' onclick = 'testfunc()' value='Testa'/></p></form><script>function testfunc() {var xhr = new XMLHttpRequest(); var fname = document.getElementById('fname').value; var lname = document.getElementById('lname').value;  var user = {}; user.firstname = fname; user.lastname = lname; xhr.open('POST', 'http://localhost:1026/user', true); xhr.setRequestHeader('Content-Type', 'application/json'); xhr.send(JSON.stringify(user))}</script></body></html>"
    t, _ := template.New("webpage").Parse(form)

    t.Execute(w, nil)
}

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

func (l *loginDB) savePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    db := l.connectToDB()
    dec := json.NewDecoder(r.Body)
    photo := Photo{}
    err := dec.Decode(&photo)
    if err != nil {
        log.Fatal(err)
    }

    res,  err := db.Prepare("insert into hat4cat.catsInHats (name, desc, image, uid) values (?, ?, ?, ?)")
    checkError(w, err)

    _, err = res.Run(photo.ImgName, photo.ImgDesc, photo.Image, photo.Uid)
    checkError(w, err)
}

func checkError(w http.ResponseWriter, err error) {
    if err != nil {
        w.WriteHeader(500) // error
        fmt.Println(err)
        fmt.Fprintf(w, "Bad input")
    }
}

func (l *loginDB) connectToDB() mysql.Conn {
    db := mysql.New("tcp", "", "130.240.170.62:3306", l.usr, l.pswd, "hat4cat")
    err := db.Connect()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Connected to database")
    return db
}

func main() {
    usr := os.Args[1]
    pswd := os.Args[2]
    login := &loginDB{usr, pswd}
    login.startWebserver()
}

