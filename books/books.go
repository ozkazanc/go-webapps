/* 
curl localhost:8080/books | jq

curl -i -X POST localhost:8080/books -d '{"title":"the hobbit","author":"jrr tolkien","year":"1937", "genre":"fantasy"}' -H "Content-Type: application/json"

GOTEST_ADMIN_PASSWD=secret go run books.go
curl localhost:8080/admin -u admin:secret
*/

package main

import (
	"fmt"
	"log"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"time"
	"strings"
	"math/rand"
	"os"
)

type Book struct {
	Id	string `json:"id"`
	Title 	string `json:"title"`
	Author	string `json:"author"`
	Year	string `json:"year"`
	Genre	string `json:"genre"`
}

type AdminPortal struct {
	password string
}

var library map[string]Book = make(map[string]Book)

const PORT = ":8080"

func main() {
	b1 := Book{Id:"0",Title:"Lotr",Author:"JRR Tolkien",Year:"1954",Genre:"Fantasy"}
	b2 := Book{Id:"1",Title:"Dune",Author:"Frank Herbert",Year:"1965",Genre:"Sci-Fi"}
	library[b1.Id] = b1
	library[b2.Id] = b2

	adminPortal := newAdminPortal()
	fmt.Println("Establishing a book database server at port", PORT)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "/books", http.StatusFound) })
	http.HandleFunc("/books", allBooks)
	http.HandleFunc("/books/", singleBook)
	http.HandleFunc("/admin", adminPortal.handler)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatal("HTTP ListenAndServe Error:", err)
	}
}

func singleBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed!"))
	}

	path := strings.Split(r.URL.Path, "/")
	if len(path) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	
	if path[2] == "random" {
		getRandomBook(w, r)	
	} else {
		getBookById(w, r, path[2])
	}
}

func getRandomBook(w http.ResponseWriter, r *http.Request) {
	ids := make([]string, len(library))
	i := 0
	for id, _ := range library {
		ids[i] = id
		i++
	}

	if i == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	} 
	
	rand.Seed(time.Now().UnixNano())
	w.Header().Add("Location", fmt.Sprintf("/books/%s", ids[rand.Intn(i)]))
	w.WriteHeader(http.StatusFound)

}

func getBookById(w http.ResponseWriter, r *http.Request, id string) {
	book, ok := library[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonBook, err := json.Marshal(book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBook)
}

func allBooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET":
			getLibrary(w, r)
		case "POST":
			postNewBook(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed!"))
	}
}

func getLibrary(w http.ResponseWriter, r *http.Request) {
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

func postNewBook(w http.ResponseWriter, r *http.Request) {
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

func newAdminPortal() *AdminPortal {
	password := os.Getenv("GOTEST_ADMIN_PASSWD")
	if password == "" {
		log.Fatal("GOTEST_ADMIN_PASSWD not set.")
	}

	return &AdminPortal{password: password}
}

func (a AdminPortal) handler(w http.ResponseWriter, r *http.Request) {
	user, passwd, ok := r.BasicAuth()
	if !ok || user != "admin" || passwd != a.password {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Unauthorized\n"))
		return
	}

	w.Write([]byte("<html><h2>Super secret admin page!!!</h2></html>"))
}


