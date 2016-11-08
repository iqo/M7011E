package main

import (
	"fmt"
	"html/template"
	"net/http"
    "log"
	//"github.com/julienschmidt/httprouter" //https://github.com/julienschmidt/httprouter

)



/*****************************************
*** Adds content on website            ***
*****************************************/
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	
    /*p := &Page{Address: dhtNode.transport.bindAddress}

    style := loadWebsite("skynet/style.html")
    script := loadWebsite("skynet/script.html")
	htmlStr := loadWebsite("skynet/webpage.html")
	t, _ := template.New("webpage").Parse(style + script + htmlStr)
    
    t.Execute(w, p)*/
    //fmt.Fprintf(w, "hello world")
    template.Must(template.ParseFiles("static/index.html")).Execute(w, nil)
}


/*****************************************
*** Starts the http-server with the    ***
*** different commands (get, post etc) ***
*****************************************/
func startWebserver() {
    //router := httprouter.New()
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
    //router.GET("/", IndexHandler)
    http.HandleFunc("/", IndexHandler)

    
    var input int
    fmt.Scan(&input)
    if input == 1 {
        fmt.Println("running on 130.240.170.62:1025")
        log.Fatal(http.ListenAndServe("130.240.170.62:1025", nil))
        } else {
            fmt.Println("running on localhost:1025")
            log.Fatal(http.ListenAndServe("localhost:1025", nil))
            
        }
}

func main() {
    startWebserver()
}

