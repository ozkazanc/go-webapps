package main

import (
	"fmt"
	"log"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"time"
)

/* 
curl localhost:8080/books | jq
curl -i -X POST localhost:8080/books -d '{"title":"hobbit","author":"jrrtolkien","year":"1937", "genre":"fantasy"}' -H "Content-Type: application/json"
*/
type Book struct {
	Id	string `json:"id"`
	Title 	string `json:"title"`
	Author	string `json:"author"`
	Year	string `json:"year"`
	Genre	string `json:"genre"`
}

var library map[string]Book = make(map[string]Book)

const PORT = ":8080"

func main() {
	b1 := Book{Id:"0",Title:"Lotr",Author:"JRR Tolkien",Year:"1954",Genre:"Fantasy"}
	b2 := Book{Id:"1",Title:"Dune",Author:"Frank Herbert",Year:"1965",Genre:"Sci-Fi"}
	library[b1.Id] = b1
	library[b2.Id] = b2

	fmt.Println("Establishing a book database server at port", PORT)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "/books", http.StatusFound) })
	http.HandleFunc("/books", books)
	//http.HandleFunc("/books/", books)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatal("HTTP ListenAndServe Error:", err)
	}
}

func books(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET":
			get(w, r)
		case "POST":
			post(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed!"))
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	lib := make([]Book, len(library))
	i := 0	
	for _, v := range library {
		lib[i] = v
		i++
	}	

	jsonLib, err := json.Marshal(lib)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonLib)
}

func post(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}


	var book Book
	err = json.Unmarshal(body, &book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	book.Id = fmt.Sprintf("%d", time.Now().UnixNano())
	library[book.Id] = book
	w.WriteHeader(http.StatusOK)
}






