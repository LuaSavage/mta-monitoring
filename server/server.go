package server

import (
   "fmt"
   "net"
   "strconv"
   "bytes"
   /*"encoding/binary"*/
)

var LENGTH_OF_INT int = 4
var LENGTH_OF_SHORT int =  2
var LENGTH_OF_CHAR int = 1
/*
var FLAGS map[string]int = map[string]int{
    "ASE_PLAYER_COUNT": 0x0004,
    "ASE_MAX_PLAYER_COUNT": 0x0008,
    "ASE_GAME_NAME": 0x0010,
    "ASE_SERVER_NAME": 0x0020,
    "ASE_GAME_MODE": 0x0040,
    "ASE_MAP_NAME": 0x0080,
    "ASE_SERVER_VER": 0x0100,
    "ASE_PASSWORDED": 0x0200,
    "ASE_SERIALS": 0x0400,
    "ASE_PLAYER_LIST": 0x0800,
    "ASE_RESPONDING": 0x1000,
    "ASE_RESTRICTION": 0x2000,
    "ASE_SEARCH_IGNORE_SECTIONS": 0x4000,
    "ASE_KEEP_FLAG": 0x8000,
    "ASE_HTTP_PORT": 0x080000,
    "ASE_SPECIAL": 0x100000,
}
*/
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

	//for {

		 s.ReadSocketData(conn)
	//}

}

func (s Server) ReadSocketData(conn *net.UDPConn) {
	fmt.Println("test this shit \n")
 	defer conn.Close() // закрываем сокет при выходе из функции

	buf := make([]byte, 1024) // буфер для чтения клиентских данных
	for {

		_, err := conn.Write([]byte("s"))

	    if err != nil {
		    fmt.Println("Write eror ", err)
		    return
	    }

	    readLen, _, err := conn.ReadFromUDP(buf) // читаем из сокета

	    if readLen > 0 {
		    if err != nil {
			    fmt.Println("ReadFromUDP eror ", err)
			    return
		    }

		    s.ReadRow(&buf)
		    //fmt.Println( string(buf))
		}
	}
}

func (s Server) ReadRow(buf *[]byte) {
	buffer := bytes.NewBuffer(*buf)

    //flags := buffer.Next(LENGTH_OF_INT)

    state:=true

    for state {

    	buffer.Next(6)
        // Length
        //len := buffer.Next(LENGTH_OF_SHORT)

        // Ip address

		for i := 0; i < LENGTH_OF_INT; i++ {
			mySlice:=buffer.Next(LENGTH_OF_CHAR)
           	fmt.Println( int(mySlice[0]))
		}

		state = false

        /*
        ip_pieces.reverse()
        server.ip = '.'.join(ip_pieces)

        server.port = buffer.read(LENGTH_OF_SHORT)
 
        if (flags & FLAGS["ASE_PLAYER_COUNT"]) != 0:
            server.playersCount = buffer.read(LENGTH_OF_SHORT)

        if (flags & FLAGS["ASE_MAX_PLAYER_COUNT"]) != 0:
            server.maxPlayersCount = buffer.read(LENGTH_OF_SHORT)

        if (flags & FLAGS["ASE_GAME_NAME"]) != 0:
            server.gameName = buffer.readString()

        if (flags & FLAGS["ASE_SERVER_NAME"]) != 0:
            server.serverName = buffer.readString()

        if (flags & FLAGS["ASE_GAME_MODE"]) != 0:
            server.modeName = buffer.readString()

        if (flags & FLAGS["ASE_MAP_NAME"]) != 0:
            server.mapName = buffer.readString()

        if (flags & FLAGS["ASE_SERVER_VER"]) != 0:
            server.verName = buffer.readString()
            
        if (flags & FLAGS["ASE_PASSWORDED"]) != 0:
            server.passworded = buffer.read(LENGTH_OF_CHAR)

        if (flags & FLAGS["ASE_SERIALS"]) != 0:
            server.serials = buffer.read(LENGTH_OF_CHAR)

        if (flags & FLAGS["ASE_PLAYER_LIST"]) != 0:
            listSize = buffer.read(LENGTH_OF_SHORT)

            for i in range(listSize):
                playerNick = buffer.readString()
                server.players.append(playerNick)
        


*/

    }
}


/*
// must be property
func (s Server) get_join_link(address string) string{
	return `mtasa://`+s.address+`:`+string(s.port)
}

*/

