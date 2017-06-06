package nim

import (
	"log"
	"net/http"
	"os"
)

// New returns a new Stack instance with no middleware preconfigured.
func New() *Stack {
	return &Stack{}
}

// Run is a convenience function that runs the nim stack as an HTTP
// server. The addr string takes the same format as http.ListenAndServe.
func Run(ns *Stack, addr ...string) {
	l := log.New(os.Stdout, "[n.] ", 0)
	address := detectAddress(addr...)
	l.Printf("Server is listening on %s", address)
	l.Fatal(http.ListenAndServe(address, ns))
}

const (
	// DefaultAddress is used if no other is specified.
	defaultServerAddress = ":3000"
)

// detectAddress
func detectAddress(addr ...string) string {
	if len(addr) > 0 {
		return addr[0]
	}
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	return defaultServerAddress
}
