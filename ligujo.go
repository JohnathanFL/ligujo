package main

import (
  "bytes"
	"database/sql"
	"fmt"
	_ "go-sqlite3"
	"log"
	"net/http"
	"strings"
	"time"
	"io"
	"encoding/json"
)


type TypeStatic struct {
  pub bool
  typeid uint
}

type NewType struct {
  name string
  statics map[string]TypeStatic
  enum *EnumType
}

/// The minimum number of parts to a path
/// Note that string.Split treats a leading `/` as an empty first element
/// So we have to +1 to all calculations
const MIN_GET_LEN = 3
const MIN_PUT_LEN = 2

type RequestServer struct{}

var db *sql.DB

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

	var parts = strings.Split(req.URL.Path, "/")
	
	if len(parts) < MIN_PUT_LEN {
  	fmt.Println("Got invalid sized request")
  	res.WriteHeader(400)
  	return
	}

  switch parts[1] {
    case "at": HandleAtUpdate(res, parts[2], req.Body)
    case "mktype": HandleMkType(res, req.Body)
    default: res.WriteHeader(400)
  }
}

func HandleAtUpdate(resp http.ResponseWriter, where string, body io.ReadCloser) {
  parts := strings.Split(where, ":")
  fmt.Printf("Handling a putat for line %v, col %v with body:\n", parts[0], parts[1])
  resp.WriteHeader(200)
}

func HandleMkType(resp http.ResponseWriter, body io.ReadCloser) {
  buf := new(bytes.Buffer)
  buf.R
  resp.WriteHeader(200)
}

// An editor has connected to us to ask for information
func HandleRequest(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Handling request...")

	var parts = strings.Split(req.URL.Path, "/")
	if len(parts) < MIN_GET_LEN {
    fmt.Println("Got invalid sized request")
    res.WriteHeader(400)
    return
	}

  switch parts[1] {
    case "at": HandleAtRequest(res, parts[2])
    default:
      fmt.Printf("Received a request for invalid method GET %v", parts[1])
      res.WriteHeader(400)
  }
}

func HandleAtRequest(resp http.ResponseWriter, where string) {
  parts := strings.Split(where, ":")
  fmt.Printf("Handling a getat for line %v, col %v\n", parts[0], parts[1])
  resp.WriteHeader(200)
  
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

	conn, err := sql.Open("sqlite3", "./type.db")
	if err != nil {
		log.Fatal("Failed to open the type.db. Did you remember to run generate.sh?")
	}
	defer conn.Close()
	db = conn

	server.ListenAndServe()
}
