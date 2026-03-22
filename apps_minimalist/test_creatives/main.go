package main

import (
	"fmt"
	"os"
	"io"
	"net/http"
	"errors"
)

// play around, to be removed
func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got / ")
	io.WriteString(w, "Welcome to server")
}

// play around, tbr
func getCoucou(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got /coucou")
	io.WriteString(w, "coucou !")
}

func getCreative(w http.ResponseWriter, r *http.Request) {
	creativeId := r.PathValue("id")
	fmt.Println("got id %s", creativeId)
	if creativeId == "test" {
		fmt.Println("serving crea")
		http.ServeFile(w, r, "./constant_campaign/banner.html")
	}
}

func main() {
	http.HandleFunc("/", getRoot)
	http.HandleFunc("GET /coucou", getCoucou)
	http.HandleFunc("GET /creatives/{id}", getCreative)

	err := http.ListenAndServe(":8090", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Print("server was closed")
	} else if err != nil {
		fmt.Print("unexpected server err: %s\n", err)
		os.Exit(1)
	}

}
