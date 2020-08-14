package main

import (
	"io/ioutil"
	"html/template"
	"net/http"
	"encoding/json"
	"flag"
	"log"
	"fmt"
	"strings"
)


type Item struct {
	Title string
	Link string
	Snippet string	
}

type Response struct {
	Items []Item
}

// -addr :PORT
var addr = flag.String("addr", ":8080", "http service address")

var templ = template.Must(template.New("sr").Parse(templateStr))

func main() {
	flag.Parse()
	log.Println("Establishing a Search Engine server at port", *addr)
    http.HandleFunc("/", search)
    err := http.ListenAndServe(*addr, nil)
    if err != nil {
        log.Fatal("HTTP ListenAndServe Error: ", err)
    }
}

func search(w http.ResponseWriter, r *http.Request) {
	// Get the URL of the search term	
	const key = "AIzaSyC7mgdbc3yXCocyHSydGNnUIByw1MIPGFY"
	const cx = "003210436456761614684:ytytx_0dwrc"
	const num = "5"

	base := "https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s&num=%s"
	query := r.FormValue("s")
	if query == "" {
		templ.Execute(w, nil)
		return
	}
	q := strings.ReplaceAll(query, " ", "+")
	
	url := fmt.Sprintf(base, key, cx, q, num)
	fmt.Println("URL is:", url)

	// Make the HTTP GET Request to Google API
	resp, err := http.Get(url)	
	if err != nil {
		log.Fatal("HTTP GET Error: ", err)	
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("IOUtil.ReadAll Error: ", err)
	}

	// Unmarshal the JSON Data
	var v Response	
	err = json.Unmarshal(body, &v)
	if err != nil {
		log.Fatal("json.Unmarshal Error: ", err)
	}
	
	/* // For debugging purposes
	fmt.Println(string(body))

	res := v.Items
	for _, val := range res {
		fmt.Println(val.Title)
	}

	fmt.Println(v.Items)
	*/

	log.Println("Searched for", query)
	templ.Execute(w, v)
}

const templateStr = `
<html>
<head>
<title>GoTest Search Engine</title>
</head>
<body>
<form action="/" name=f method="GET">
    <input maxLength=1024 size=70 name=s value="" title="Text to search">
    <input type=submit value="Search" name=search>
</form>

<div class="results">
	{{range .Items}}
	<div class="item">
		<h3>{{.Title}}</h3>
		<a href="{{.Link}}" target="_blank">{{.Link}}</a>
		<p>{{.Snippet}}</p>
		<br>
	</div>
	{{end}}
</div>
</body>
</html>
`
