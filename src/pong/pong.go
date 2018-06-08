package main

import (
	"fmt"
	"net/http"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Healthy")
}

func main() {
	http.HandleFunc("/healthz", Healthz)
	http.ListenAndServe(":8080", nil)
}
