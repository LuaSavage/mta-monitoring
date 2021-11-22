# MTA Server Monitoring

Light weighted solution for MTA Server monitoring from ASE port.
A.S.E means All Seeing Eye and actually its udp port

This library is adoptation of https://github.com/Lipau3n/mtasa-monitoring to golang

## Usage
```go
import (
	"mta-monitoring/server"
)

func main() {
  	// pass server address and port
	newServer :=server.NewServer("217.106.106.107", 22044)
  
  	// print link to foin mta sa server
	fmt.Println(newServer.Get_join_link())
}
```

## Server information
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
