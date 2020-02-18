package main

import (
	"fmt"
	"net/http"
	"time"
)

type RequestServer struct{}

func (r RequestServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
    // We'll allow PUT/POST as synonyms since I can't find a good reason to restrict.
    if req.Method == "PUT" || req.Method == "POST" {
        HandleUpdate(resp, req)
    } else if req.Method == "GET" {
        HandleRequest(resp, req)
    }
}

// An interpreter has connected to us to give us new information
func HandleUpdate(res http.ResponseWriter, req *http.Request) {
    fmt.Println("Handling update...")
}

// An editor has connected to us to ask for information
func HandleRequest(res http.ResponseWriter, req *http.Request) {
    fmt.Println("Handling request...")
}

func main() {
	fmt.Println("Saluton, mondo!")
	server := &http.Server{
		Addr:           ":8080",
		Handler:        &RequestServer{},
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	server.ListenAndServe()
}
