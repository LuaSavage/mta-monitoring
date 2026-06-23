![test workflow](https://github.com/LuaSavage/mta-monitoring/actions/workflows/go.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/LuaSavage/mta-monitoring)](https://goreportcard.com/report/github.com/LuaSavage/mta-monitoring)
[![codecov](https://codecov.io/gh/LuaSavage/mta-monitoring/branch/main/graph/badge.svg?token=FUPH9E2C38)](https://codecov.io/gh/LuaSavage/mta-monitoring)

# MTA Server Monitoring

Lightweight library for monitoring MTA:SA servers via the ASE (All-Seeing Eye) UDP port.

Inspired by https://github.com/Lipau3n/mtasa-monitoring

The library runtime uses only the Go standard library. Test dependencies are `testify` and `go.uber.org/mock`.

## Getting started

### Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/LuaSavage/mta-monitoring/server"
)

func main() {
	exampleServer := server.NewServer("185.71.66.81", 22003)

	// UpdateOnce fetches ASE data once. Use a poller for periodic updates.
	if err := exampleServer.UpdateOnce(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", exampleServer)
	fmt.Println(exampleServer.GetJoinLink())
}
```

### Output

```shell
&{Timeout:0 Game:mta Address:185.71.66.81 Port:22003 AsePort:22126 Name:Actual-server-name Gamemode:RPG Map:None Version:1.5n Somewhat:0 Players:280 Maxplayers:815}
mtasa://185.71.66.81:22003
```

### Timeout

Set `Server.Timeout` to the UDP read/write deadline in seconds. When zero, a 5 second default is used.

```go
exampleServer := server.NewServer("185.71.66.81", 22003)
exampleServer.Timeout = 10
```

### Server information

* **Game** — game identifier (`mta`)
* **Address** — MTA server IP address
* **Port** — main server port (UDP)
* **AsePort** — ASE port (`main port + 123`)
* **Name** — server name
* **Gamemode** — server mode
* **Map** — current map
* **Version** — MTA:SA server version
* **Somewhat** — ASE passworded flag
* **Players** — current player count
* **Maxplayers** — maximum player slots

## Build

Run tests:

```shell
make test
```

Regenerate mocks:

```shell
make generate
```

Build the example CLI:

```shell
make build-example
```

Build and run the example (optional address and port arguments):

```shell
make run-example
go run ./cmd/example 185.71.66.81 22003
```
