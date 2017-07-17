package test

import (
	"io"
    "net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	payload, _ := os.Open("./test/tickets.json")
	io.Copy(w, payload)
}

func MockZD() {
	http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
