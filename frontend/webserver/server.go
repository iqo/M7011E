package main

import (
	"fmt"
	//"html/template"
	"net/http"
    "log"
	"github.com/julienschmidt/httprouter" //https://github.com/julienschmidt/httprouter

)



/*****************************************
*** Adds content on website            ***
*****************************************/
func IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	
    /*p := &Page{Address: dhtNode.transport.bindAddress}

    style := loadWebsite("skynet/style.html")
    script := loadWebsite("skynet/script.html")
	htmlStr := loadWebsite("skynet/webpage.html")
	t, _ := template.New("webpage").Parse(style + script + htmlStr)

    t.Execute(w, p)*/
    fmt.Fprintf(w, "hello world")
}


/*****************************************
*** Starts the http-server with the    ***
*** different commands (get, post etc) ***
*****************************************/
func startWebserver() {
    router := httprouter.New()
    router.GET("/", IndexHandler)
    log.Fatal(http.ListenAndServe("localhost:1025", router))
}

