// Package nim is a lightweight middleware stack manager in Golang
// that is based on idiomatic net/http method signatures.
//
// For a full guide visit http://github.com/nimgo/nim
//
//  package main
//
//  import (
//    "net/http"
//    "fmt"
//    "github.com/nimgo/nim"
//  )
//
//  func main() {
//    mux := http.NewServeMux()
//    mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
//      fmt.Fprintf(w, "Welcome to nimble!")
//    })
//
//    n := nimble.Default()
//    n.Use(mux)
//    n.Run(":8000")
//  }
package nim
