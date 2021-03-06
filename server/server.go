package server

import (
	"bytes"
	"fmt"
	"net"
	"reflect"
	"strconv"
)

type Server struct {
	Timeout    float64
	Game       string
	Address    string
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
	newServer := Server{
		Timeout:    0,
		Game:       "",
		Address:    address,
		Port:       port,
		AsePort:    port + 123,
		Name:       "",
		Gamemode:   "",
		Map:        "",
		Version:    "",
		Somewhat:   "",
		Players:    0,
		Maxplayers: 0,
	}
	_, connection := newServer.Connect()
	_, responseBuffer := newServer.ReadSocketData(connection)
	newServer.ReadRow(responseBuffer)
	return &newServer
}

func (s *Server) Connect() (*Server, *net.UDPConn) {

	updAddr, err := net.ResolveUDPAddr("udp", s.Address+":"+strconv.Itoa(s.AsePort))

	if err != nil {
		fmt.Println("ResolveUDPAddr failed", err)
		return s, nil
	}

	conn, err := net.DialUDP("udp", nil, updAddr)

	if err != nil {
		fmt.Println("Couldn't establish UDP connection. \n", err)
		return s, conn
	}

	return s, conn
}

func (s *Server) ReadSocketData(conn *net.UDPConn) (*Server, *[]byte) {

	defer conn.Close() // закрываем сокет при выходе из функции

	buf := make([]byte, 1024) // буфер для чтения клиентских данных

	for {

		_, err := conn.Write([]byte("s"))

		if err != nil {
			fmt.Println("Write eror ", err)
			return s, &buf
		}

		readLen, _, err := conn.ReadFromUDP(buf) // читаем из сокета

		if readLen > 0 {
			if err != nil {
				fmt.Println("ReadFromUDP eror ", err)
				return s, &buf
			}

			break
		}

	}

	return s, &buf

}

func (s *Server) ReadRow(buf *[]byte) *Server {

	buffer := bytes.NewBuffer(*buf)

	params := [9]string{"Game", "Port", "Name", "Gamemode", "Map", "Version", "Somewhat", "Players", "Maxplayers"}

	//reading begins from 4 byte
	buffer.Next(4)

	obj := reflect.Indirect(reflect.ValueOf(s))

	for i := 0; i < len(params); i++ {

		length := int(buffer.Next(1)[0]) - 1
		value := buffer.Next(length)
		fieldName := params[i]
		field := obj.FieldByName(fieldName)

		if field.Type().Name() == "int" {
			i, _ := strconv.Atoi(string(value))
			field.SetInt(int64(i))
		}

		if field.Type().Name() == "string" {
			field.SetString(string(value))
		}

	}

	return s
}

func (s Server) GetJoinLink() string {
	// return link to join mta sa server
	return `mtasa://` + s.Address + `:` + strconv.Itoa(s.Port)
}
