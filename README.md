![test workflow](https://github.com/LuaSavage/mta-monitoring/actions/workflows/go.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/LuaSavage/mta-monitoring)](https://goreportcard.com/report/github.com/LuaSavage/mta-monitoring)
[![codecov](https://codecov.io/gh/LuaSavage/mta-monitoring/branch/master/graph/badge.svg?token=7TMGVKVMET)](https://codecov.io/gh/LuaSavage/mta-monitoring)

# MTA Server Monitoring

Light weighted solution for MTA Server monitoring from ASE port.  
A.S.E means All Seeing Eye and actually its udp port

Inspired by https://github.com/Lipau3n/mtasa-monitoring  
Depends only on standard libraries.

## Getting started
### Example
```go
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
```

### Server information
* **Game** (mta)
* **Address** string with MTA server ip address
* **Port** - server main port (UDP)
* **AsePort** - main MTA:SA port + 123
* **Name** - server name
* **Gamemode** - server mode
* **Map** - server map
* **Version** - mta:sa server version
* **Players** - number of players on the server right now
* **Maxplayers** - the maximum number of players that can join

## Build 
You can modify example and then build it:
```shell
make build-example
```
Or build and run simultaneously:
```shell
make run-example
```