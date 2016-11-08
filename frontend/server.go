package main

import (
	//"fmt"
	"html/template"
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
    //fmt.Fprintf(w, "hello world")
    template.Must(template.ParseFiles("static/index.html")).Execute(w, nil)
}


/*****************************************
*** Starts the http-server with the    ***
*** different commands (get, post etc) ***
*****************************************/
func startWebserver() {
    router := httprouter.New()
    //http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
    router.GET("/", IndexHandler)
    router.ServeFiles("/*filepath", http.Dir("static"))
    /*
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))
    http.Handle("/css/", fs)
    http.Handle("/js/", fs)
    http.Handle("/bootstrap/", fs)
    http.Handle("/fonts/", fs)
    http.Handle("/img/", fs)
*/
    log.Fatal(http.ListenAndServe("130.240.170.62:1025", router))
}

func main() {
    startWebserver()
}

