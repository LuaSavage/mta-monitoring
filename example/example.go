package main

import (
	"fmt"

	"github.com/LuaSavage/mta-monitoring/server"
)

func main() {
	// pass server address and port
	exampleServer := server.NewServer("185.71.66.81", 22003)

	// Note that it updating fields once.
	// To update them frequently or on occasion you've to have some sort of poller
	if err := exampleServer.UpdateOnce(); err != nil {
		panic(err)
	}

	// Printing updated data in objects structure
	fmt.Printf("%+v\n", exampleServer)

	// Printing link to join mta:sa server
	fmt.Println(exampleServer.GetJoinLink())
}
