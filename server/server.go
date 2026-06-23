package server

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"reflect"
	"strconv"
	"time"
)

//go:generate go run go.uber.org/mock/mockgen@latest -source ./server.go -destination=./udpconn_mock.go -package=server

const (
	defaultTimeout   = 5 * time.Second
	maxReadAttempts  = 3
	aseQueryPayload  = "s"
	aseResponseBufSz = 1024
)

// UDPconnection is the minimal UDP surface used for ASE queries and testing.
type UDPconnection interface {
	Write(b []byte) (n int, err error)
	ReadFromUDP(b []byte) (int, *net.UDPAddr, error)
	Close() error
}

// Server holds MTA:SA server connection settings and fields populated from an ASE response.
type Server struct {
	// Timeout is the UDP read/write deadline in seconds. Zero uses a 5 second default.
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

// NewServer creates a Server for the given game address and main UDP port.
// The ASE port is computed as port + 123.
func NewServer(address string, port int) *Server {
	return &Server{
		Address: address,
		Port:    port,
		AsePort: port + 123,
	}
}

func (s *Server) timeout() time.Duration {
	if s.Timeout == 0 {
		return defaultTimeout
	}
	return time.Duration(s.Timeout * float64(time.Second))
}

func (s *Server) refreshDeadline() error {
	conn, ok := s.connection.(*net.UDPConn)
	if !ok {
		return nil
	}

	deadline := time.Now().Add(s.timeout())
	if err := conn.SetReadDeadline(deadline); err != nil {
		return fmt.Errorf("set read deadline: %w", err)
	}
	if err := conn.SetWriteDeadline(deadline); err != nil {
		return fmt.Errorf("set write deadline: %w", err)
	}
	return nil
}

// Connect opens a UDP connection to the server's ASE port.
func (s *Server) Connect() (*net.UDPConn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", s.Address+":"+strconv.Itoa(s.AsePort))
	if err != nil {
		return nil, fmt.Errorf("resolve UDPAddr failed: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, fmt.Errorf("couldn't establish udp connection: %w", err)
	}

	s.connection = conn
	if err := s.refreshDeadline(); err != nil {
		_ = conn.Close()
		s.connection = nil
		return nil, err
	}

	return conn, nil
}

// ReadSocketData sends the ASE status query and reads the response into a buffer.
func (s *Server) ReadSocketData() (*[]byte, error) {
	buf := make([]byte, aseResponseBufSz)

	for attempt := 0; attempt < maxReadAttempts; attempt++ {
		if err := s.refreshDeadline(); err != nil {
			return nil, err
		}

		if _, err := s.connection.Write([]byte(aseQueryPayload)); err != nil {
			return nil, err
		}

		readLen, _, err := s.connection.ReadFromUDP(buf)
		if err != nil {
			if isTimeoutError(err) {
				return nil, fmt.Errorf("read ASE response timed out: %w", err)
			}
			return nil, err
		}

		if readLen > 0 {
			return &buf, nil
		}
	}

	return nil, errors.New("empty ASE response after retries")
}

func isTimeoutError(err error) bool {
	if errors.Is(err, os.ErrDeadlineExceeded) {
		return true
	}

	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

// ReadRow parses an ASE response buffer and applies the fields to the server.
func (s *Server) ReadRow(buf *[]byte) *Server {
	buffer := bytes.NewBuffer(*buf)
	params := [9]string{"Game", "Port", "Name", "Gamemode", "Map", "Version", "Somewhat", "Players", "Maxplayers"}

	// reading begins from 4 byte
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

// UpdateOnce connects if needed, queries ASE once, parses the response, and closes the connection.
func (s *Server) UpdateOnce() error {
	if s.connection == nil {
		if _, err := s.Connect(); err != nil {
			return err
		}
	}

	defer func() {
		s.connection.Close()
		s.connection = nil
	}()

	buff, err := s.ReadSocketData()
	if err != nil {
		return err
	}

	s.ReadRow(buff)
	return nil
}

// GetJoinLink returns the mtasa:// join link shown in the in-game browser.
func (s *Server) GetJoinLink() string {
	return fmt.Sprintf("mtasa://%s:%d", s.Address, s.Port)
}
