package server

import (
   "fmt"
   "net"
)

type Server struct {
	timeout    float64
	game       string
	address	   string
	port       int
	asePort    int
	name       string
	gamemode   string
	map_name   string
	version    string
	somewhat   string
	players    int
	maxplayers int
}

func NewServer(address string, port int) *Server {
	newServer:= Server{address: address, port: port, asePort: port + 123}
	newServer.Connect()
 	return &newServer
}

func (s Server) Connect (){

	fmt.Println(" port ", s.asePort)
	//updAddr :=net.UDPAddr{IP: net.ParseIP(s.address), Port: s.asePort }
	updAddr, err := net.ResolveUDPAddr("udp", string(s.address)+":"+string(s.asePort))

	if err != nil {
		fmt.Println(" ResolveUDPAddr failed", err)
		return
	}

	//listener, err := net.ListenUDP("udp", &updAddr)
	conn, err := net.DialUDP("udp", nil, updAddr)

	if err != nil {
		fmt.Println("Could not establish UDP connection. \n", err)
		return
	}

	fmt.Println(updAddr,conn)

	for {

		/*go*/ //s.ReadSocketData(conn)
	}

}

func (s Server) ReadSocketData(conn *net.UDPConn) {
 	defer conn.Close() // закрываем сокет при выходе из функции


	buf := make([]byte, 32) // буфер для чтения клиентских данных
	for {

	      readLen, _, err := conn.ReadFromUDP(buf) // читаем из сокета
	      if err != nil {
		      fmt.Println(err)
		      return
	      }
	      fmt.Println(readLen)

	}
}


/*
// must be property
func (s Server) get_join_link(address string) string{
	return `mtasa://`+s.address+`:`+string(s.port)
}

*/

