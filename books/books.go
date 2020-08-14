package main

import (
	"fmt"
	"log"
	//"encoding/json"
	"net/http"
)

const PORT = ":8080"

func main() {
	fmt.Println("Establishing a book database server at port", PORT)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "/books", http.StatusFound) })
	http.HandleFunc("/books", books)
	http.HandleFunc("/books/", books)
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
	fmt.Fprintf(w, "In get func\n")
}

func post(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "In post func\n")
}
