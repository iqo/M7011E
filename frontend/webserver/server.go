package dht

import (
	"fmt"
	//"html/template"
	"net/http"
	"github.com/julienschmidt/httprouter" //https://github.com/julienschmidt/httprouter

)

type Page struct {
    Address     string
}


/*****************************************
*** Adds content on website            ***
*****************************************/
func (dhtNode *DHTNode) IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	
    /*p := &Page{Address: dhtNode.transport.bindAddress}

    style := loadWebsite("skynet/style.html")
    script := loadWebsite("skynet/script.html")
	htmlStr := loadWebsite("skynet/webpage.html")
	t, _ := template.New("webpage").Parse(style + script + htmlStr)

    t.Execute(w, p)*/
    fmt.Fprint("hello world", w)
}

/*****************************************
*** Opens a file and reads the content ***
*** one line at the time, the whole    ***
*** file is returned as one string     ***
*****************************************/
func loadWebsite(filename string) string{
	htmlStr := ""
    file, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        htmlStr = htmlStr + scanner.Text() + "\n"
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
    return htmlStr
}

/*****************************************
*** Starts the http-server with the    ***
*** different commands (get, post etc) ***
*****************************************/
func (dhtNode *DHTNode) startWebserver() {
    router := httprouter.New()
    router.GET("/", dhtNode.IndexHandler)
    log.Fatal(http.ListenAndServe(dhtNode.contact.ip+":"+dhtNode.contact.port, router))
}

