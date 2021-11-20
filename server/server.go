package server

import (
   "fmt"
   "net"
   "strconv"
   "bytes"
   "reflect"
)

type Server struct {
	Timeout    float64
	Game       string
	Address	   string
	Port       int
	AsePort    int
	Name       string
	Gamemode   string
	Map        string
	Version    string
	Somewhat   string
	Players    int
	Maxplayers int
}

func NewServer(address string, port int) *Server {
	newServer:= Server{Address: address, Port: port, AsePort: port + 123}
	newServer.Connect()
 	return &newServer
}

func (s Server) Connect (){

	updAddr, err := net.ResolveUDPAddr("udp", s.Address+":"+strconv.Itoa(s.AsePort))

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

    params:= [9]string{"Game", "Port", "Name", "Gamemode", "Map", "Version", "Somewhat", "Players", "Maxplayers"}
 
    //reading begins from 4 byte
    buffer.Next(4)

    for i:=0; i<len(params); i++ {
        length := int(buffer.Next(1)[0])-1
        value := buffer.Next(length)
        fmt.Println( string(value))    


        fieldName:=params[i]

        obj := reflect.Indirect(reflect.ValueOf(&s))
        field:=obj.FieldByName(fieldName)
        fmt.Println(field,fieldName)
        if field.Type().Name() == "int"{
            i, _ := strconv.Atoi(string(value))
            field.SetInt(int64(i))
        } 

        if field.Type().Name() == "string"{
            field.SetString(string(value))
        }          

    } 

    fmt.Println(s)     
}


/*
// must be property
func (s Server) get_join_link(address string) string{
	return `mtasa://`+s.address+`:`+string(s.port)
}

*/

