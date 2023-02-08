// Package server implements mta:sa server description object with ability to self update
package server

import (
	"bytes"
	"fmt"
	"net"
	"reflect"
	"strconv"
)

//go:generate mockgen -source ./server.go -destination=./../mock/udpconn_mock.go -package=mock
type UDPconnection interface {
	Write(b []byte) (n int, err error)
	ReadFromUDP(b []byte) (int, *net.UDPAddr, error)
	Close() error
}

/*
 This is a server object.

 All exported fields contains game server metadata
*/
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
	connection UDPconnection
}

// Constructor of mta:sa server object.
// Depends on mta:sa server address and port.
// The address and port are the same as what you may typed in the game browser address row
func NewServer(address string, port int) (server *Server) {
	server = &Server{
		Address: address,
		Port:    port,
		AsePort: port + 123,
	}

	return
}

// It establishing udp connection to the game server
func (s *Server) Connect() (conn *net.UDPConn, err error) {
	updAddr, err := net.ResolveUDPAddr("udp", s.Address+":"+strconv.Itoa(s.AsePort))
	if err != nil {
		err = fmt.Errorf("resolve UDPAddr failed: %s", err)
		return
	}

	conn, err = net.DialUDP("udp", nil, updAddr)
	if err != nil {
		err = fmt.Errorf("couldn't establish udp connection: %s", err)
		return
	}
	s.connection = conn

	return
}

// Reading data to buffer from socket
func (s *Server) ReadSocketData() (*[]byte, error) {
	buf := make([]byte, 1024)

	for {
		_, err := s.connection.Write([]byte("s"))
		if err != nil {
			return nil, err
		}

		readLen, _, err := s.connection.ReadFromUDP(buf)
		if err != nil {
			return nil, err
		}

		if readLen > 0 {
			break
		}
	}

	return &buf, nil
}

// Interpreting data from the buffer and applying it to the server object field by field
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

// Establishing connection, reading, interpreting data and wraps up
func (s *Server) UpdateOnce() (err error) {
	if s.connection == nil {
		if _, err = s.Connect(); err != nil {
			return err
		}
	}

	defer func() {
		s.connection.Close()
		s.connection = nil
	}()

	buff, err := s.ReadSocketData()
	if err != nil {
		return
	}

	s.ReadRow(buff)
	return
}

// It returns string with join link,
// which maybe you got used to see in ingame browser.
func (s *Server) GetJoinLink() string {
	// return link to join mta sa server
	return fmt.Sprintf("mtasa://%s:%d", s.Address, s.Port)
}
