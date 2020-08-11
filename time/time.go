package main

import (
	"fmt"
	"log"
	"time"
	"net/http"
)

func main() {
	fmt.Printf("Time is now! %s", time.Now())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "/time", http.StatusFound)	})

	http.HandleFunc("/time", tellTime)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Test Error!!!")
	}
}

func tellTime(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Time now is...\n%s", time.Now())
}
