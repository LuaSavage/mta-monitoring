package main

import (
	"fmt"
	"mta-monitoring/server"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	newServer := server.NewServer("217.106.106.107", 22044)
	spew.Dump(newServer)
	fmt.Println(newServer.GetJoinLink())
}
