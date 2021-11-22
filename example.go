package main

import (
	"github.com/davecgh/go-spew/spew"
	"mta-monitoring/server"
)

func main() {
	newServer := server.NewServer("217.106.106.107", 22044)
	spew.Dump(newServer)
}
