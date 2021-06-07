package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if name := r.FormValue("name"); name != "" {
		fmt.Fprintf(w, "Hello, %s!", name)
	} else {
		fmt.Fprintf(w, "Hello, World!!")
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
