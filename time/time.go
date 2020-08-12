package main

import (
	"fmt"
	"log"
	"time"
	"net/http"
)

func main() {
	fmt.Println("Establishing a time server at port 8080.")
	fmt.Printf("Time at the start is %s\n", time.Now())
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "/time", http.StatusFound) })
	http.HandleFunc("/time", tellTime)
	
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("HTTP Server Error: ", err)
	}
}

func tellTime(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Time now is...\n%s", time.Now())
}
