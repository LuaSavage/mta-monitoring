package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/LuaSavage/mta-monitoring/server"
)

const defaultAddress = "185.71.66.81"
const defaultPort = 22003

func main() {
	address := defaultAddress
	port := defaultPort

	if len(os.Args) > 1 {
		address = os.Args[1]
	}
	if len(os.Args) > 2 {
		parsedPort, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("invalid port %q: %v", os.Args[2], err)
		}
		port = parsedPort
	}

	exampleServer := server.NewServer(address, port)

	// UpdateOnce fetches ASE data once. Use a poller for periodic updates.
	if err := exampleServer.UpdateOnce(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", exampleServer)
	fmt.Println(exampleServer.GetJoinLink())
}
