package main

import (
  //"bytes"
  "errors"
	"database/sql"
	"fmt"
	_ "go-sqlite3"
	"log"
	"net/http"
	"strings"
	"time"
	"io"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

type TypeID uint64

func (id TypeID) IsPrimitive() bool { return id <= 2052 }

// Should be stored in a []Field or 
type Field struct {
  Access int32 `json:"access"` // One of 0,1,2 for priv, readonly, or pub
  Name string `json:"name"`
  Contains TypeID `json:"typeid"` // The type we hold
}

type Type struct {
  Name string `json:"name"`
  Statics []Field `json:"statics"`
  Pos string `json:"pos"` // line:col
  // One of the TypeIDs for different types of types
  // If it's <= 2052, it's a primitive. Otherwise, it uses that operator
  Ty TypeID `json:"ty"`
  // If it's an enum, this is the tagType
  // Otherwise it's currently unused
  // Called `Contained` in SQL
  Backing TypeID `json:"backing"`
  // Only used for Arrays and Tuples
  Len uint
  // If it's an enum, then pub is always true (ignored)
  // If it's a tuple, they're numbered 0..=n
  Fields []Field `json:"fields"`
}

func GetTypeFromID(id TypeID) (res *Type, err error) {
  fmt.Printf("Retrieving #%v\n", id)
  row := db.QueryRow("select name, pos, type, contained from Types where id = ?;", id)
  if row == nil { res = nil; err = errors.New("No type row found.") ; return }

  res = new(Type)
  
  err = row.Scan(&res.Name, &res.Pos, &res.Ty, &res.Backing)
  if err != nil { res = nil; err = errors.New("No type row found."); return }
  fmt.Printf("  Found type %v @ %v of ty %v backed by %v\n", res.Name, res.Pos, res.Ty, res.Backing)

  {
    allStatics, _ := db.Query("select name, access, contains from Fields where type = ? and static = 1;", id)

    res.Statics = make([]Field, 0)
    for allStatics.Next() {
      var static Field
      allStatics.Scan(&static.Name, &static.Access, &static.Contains)
      res.Statics = append(res.Statics, static)
    }
  }
  {
    allFields, _ := db.Query("select name, access, contains from Fields where type = ? and static = 0;", id)

    res.Statics = make([]Field, 0)
    for allFields.Next() {
      var field Field
      err := allFields.Scan(&field.Name, &field.Access, &field.Contains)
      fmt.Printf("    Err was %v\n", err)
      fmt.Printf("    Found field %#v with access %v containing %v\n", field.Name, field.Access, field.Contains)
      res.Fields = append(res.Fields, field)
    }
  }

  err = nil
  return 
}

/// The minimum number of parts to a path
/// Note that string.Split treats a leading `/` as an empty first element
/// So we have to +1 to all calculations
const MIN_GET_LEN = 3
const MIN_PUT_LEN = 2

type RequestServer struct{}

var db *sql.DB
var insertStmt *sql.Stmt

func HasType(id TypeID) bool {
  if db == nil { return false }

  rows, _ := db.Query("select COUNT(*) from Types where id = ?;", id)
  var count uint32
  rows.Scan(&count)
  return count != 0
}

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
    default: res.WriteHeader(500)
  }
}

func HandleAtUpdate(resp http.ResponseWriter, where string, body io.ReadCloser) {
  parts := strings.Split(where, ":")
  fmt.Printf("Handling a putat for line %v, col %v with body:\n", parts[0], parts[1])
  resp.WriteHeader(200)
}

func HandleMkType(resp http.ResponseWriter, body io.ReadCloser) {
  fmt.Printf("Doing /mktype\n")
  b, _ := ioutil.ReadAll(body)
  
  if len(b) == 0 {
    // TODO: Get some better error codes in here
    resp.WriteHeader(403)
    return
  }
  //fmt.Printf("Body was:\n%s\n", b)

  // Map of typeid -> type
  var types map[TypeID]Type
  err := json.Unmarshal(b, &types)
  if err != nil {
    fmt.Printf("  Err was %v\n", err)
    resp.WriteHeader(400)
    return
  }

  
  for id, ty := range types {
    // Already know about it, and this is mktype,
    // not some `updatetype`
    if HasType(id) { continue }

    fmt.Printf("Adding #%v(%v)", id, ty)

    _, err := db.Exec(
      "insert into Types (id, name, pos, type, contained, len) VALUES (?, ?, ?, ?, ?, ?);",
      id, ty.Name, ty.Pos, ty.Ty, ty.Backing, ty.Len,
    )
    if err != nil {
      fmt.Printf("Insert err was: %v", err)
      resp.WriteHeader(400)
      return
    }

    for _, static := range ty.Statics {
      _, err := db.Exec(
        "insert into Fields (static, type, name, access, contains, isRet) VALUES (1, ?, ?, ?, ?, NULL);",
        id, static.Name, static.Access, static.Contains,
      )
      if err != nil {
        fmt.Printf("Insert static err was: %v", err)
        resp.WriteHeader(400)
        return
      }
    }

    for _, field := range ty.Fields {
      _, err := db.Exec(
        "insert into Fields (static, type, name, access, contains, isRet) VALUES (0, ?, ?, ?, ?, NULL);",
        id, field.Name, field.Access, field.Contains,
      )

      if err != nil {
        fmt.Printf("Insert field err was: %v", err)
        resp.WriteHeader(400)
        return
      }
    }
    
  }
  
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
    case "typeid": HandleTypeIDRequest(res, parts[2])
    default:
      fmt.Printf("Received a request for invalid method GET %v", parts[1])
      res.WriteHeader(400)
  }
}

func HandleTypeIDRequest(resp http.ResponseWriter, id string) {
  realID, err := strconv.ParseUint(id, 10, 32)
  if err != nil { resp.WriteHeader(400); return }
  ty, err := GetTypeFromID(TypeID(realID))
  if err != nil { resp.WriteHeader(400); return }

  js, err := json.Marshal(ty)
  if err != nil { resp.WriteHeader(500); return }

  resp.WriteHeader(200)
  resp.Write(js)
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
	insertStmt, _ = db.Prepare("insert into Types (id, name, pos, type, contained, len) VALUES (?, ?, ?, ?, ?, ?);")
	defer insertStmt.Close()

	server.ListenAndServe()
}
