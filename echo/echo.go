package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Establishing an echo server at port 8080.")
	http.HandleFunc("/", echo)
	log.Fatal("HTTP Server Error: ", http.ListenAndServe(":8080", nil))
}

func echo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", r.URL.Path[1:])
}
