package server

import (
   "fmt"
   "net"
   "strconv"
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

	updAddr, err := net.ResolveUDPAddr("udp", s.address+":"+strconv.Itoa(s.asePort))

	if err != nil {
		fmt.Println(" ResolveUDPAddr failed", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, updAddr)

	if err != nil {
		fmt.Println("Could not establish UDP connection. \n", err)
		return
	}
	fmt.Println(conn)
	/*for {

		go s.ReadSocketData(conn)
	}*/

}

func (s Server) ReadSocketData(conn *net.UDPConn) {
	fmt.Println("test this shit \n")
 	defer conn.Close() // закрываем сокет при выходе из функции


	buf := make([]byte, 1024) // буфер для чтения клиентских данных
	for {

	      readLen, _, err := conn.ReadFromUDP(buf) // читаем из сокета
	      if err != nil {
		      fmt.Println("ReadFromUDP eror ", err)
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

