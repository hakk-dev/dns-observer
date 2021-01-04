package main 

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type Tracker struct {
	Id 		string
}

var (
	http_addr string
	http_port int
)

var templateStr = `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<title>Which DNS Provider?</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="stylesheet" type="text/css" href="//{{.Id}}.dns-observer.e42.xyz/css/main.css">
</head>
<body>
	<h1>Which DNS Provider?</h1>
	<p>A neat little program to determine which DNS provider a visitor is using.</p>
	<p id="provider"></p>
	<script type="text/javascript">
		var url = "/api";
		var params = "key={{.Id}}";
		var http = new XMLHttpRequest();

		http.open("GET", url+"?"+params, true);
		http.onreadystatechange = function()
		{
		    if(http.readyState == 4 && http.status == 200) {
		        var obj = JSON.parse(http.responseText);
		        document.getElementById("provider").innerHTML = "Device DNS provider is: " + obj.Name;
		    }
		}
		http.send(null);
	</script>
</body>
</html>`

func main() {
	// load the config
	configSetup()

	// start the dns server
	dnsStart()

	// start the http server
	httpStart()
}

func httpStart() {
	http.HandleFunc("/", templatedHandler)
	http.HandleFunc("/api", apiHandler)
	log.Printf("Starting on port %s:%d", http_addr, http_port)
	http.ListenAndServe(fmt.Sprintf("%s:%d", http_addr, http_port), nil)
}

func templatedHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Header)
	log.Println(r.Header.Get("X-Forwarded-For"))
	tmplt := template.New("hello world")
	tmplt, _ = tmplt.Parse(templateStr)

	uid, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}

	p := Tracker{Id: uid.String()}

	tmplt.Execute(w, p)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	var demo Provider

	keys, ok := r.URL.Query()["key"]
    
    if !ok || len(keys[0]) < 1 {
        log.Println("key is missing")
        demo = Provider{"key is missing"}
    } else {
    	demo = db[keys[0]]
    }

	jsonData, err := json.Marshal(demo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
