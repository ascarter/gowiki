package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alpental/gowiki"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", r.URL.Path[1:])
}

func main() {
	// Register handlers
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", gowiki.ViewHandler)
	http.HandleFunc("/edit/", gowiki.EditHandler)
	http.HandleFunc("/save/", gowiki.SaveHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
