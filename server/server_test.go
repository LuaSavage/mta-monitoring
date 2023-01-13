package server

import (
	"bytes"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const typicalResponse = `"EYE1\x04mta\x0622044M                          MTA:SA Türkiye - Norm Gaming [ Turkish / Turkey ]\aMTA:SA\x05None\x041.5\x021\x020\x0370\x01"`
const testIp = "217.106.106.107"
const testPort = 22044

func GetTypicalBytes(t *testing.T) (buf []byte) {
	t.Helper()

	unquotedTypicalResponse, _ := strconv.Unquote(typicalResponse)
	buf = []byte(unquotedTypicalResponse)
	return
}

func ValidateFields(t *testing.T, testServer *Server) {
	t.Helper()

	expectedValues := make(map[string]interface{})
	expectedValues["Game"] = "mta"
	expectedValues["Port"] = 22044
	expectedValues["Name"] = "                          MTA:SA Türkiye - Norm Gaming [ Turkish / Turkey ]"
	expectedValues["Gamemode"] = "MTA:SA"
	expectedValues["Map"] = "None"
	expectedValues["Version"] = "1.5"
	expectedValues["Somewhat"] = "1"
	expectedValues["Players"] = 0
	expectedValues["Maxplayers"] = 70

	testServerObj := reflect.Indirect(reflect.ValueOf(testServer))

	for key, value := range expectedValues {
		field := testServerObj.FieldByName(key)
		assert.Equal(t, field.Interface(), value, "Field "+key+" must be equal")
	}
}

func TestReadSocketData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUdpConn := NewMockUDPconnection(mockCtrl)
	mockUdpConn.EXPECT().Write(gomock.AssignableToTypeOf([]byte("s"))).Return(1, nil)

	emptyByte := []byte("")

	bytesOfTypicalResponse := GetTypicalBytes(t)
	mockUdpConn.EXPECT().ReadFromUDP(gomock.AssignableToTypeOf(emptyByte)).DoAndReturn(func(b []byte) (_ int, _ *net.UDPAddr, _ error) {
		copy(b, bytesOfTypicalResponse)
		return len(bytesOfTypicalResponse), nil, nil
	}).Times(1)

	newServer := NewServer( /* address */ testIp /* port */, testPort)
	newServer.connection = mockUdpConn
	recievedBytes, err := newServer.ReadSocketData()

	assert.NoError(t, err)
	recievedBytesShort := (*recievedBytes)[:len(bytesOfTypicalResponse)]
	assert.True(t, bytes.Equal(bytesOfTypicalResponse, recievedBytesShort), "received bytes unequal")
}

func TestReadRow(t *testing.T) {
	bytesOfTypicalResponse := GetTypicalBytes(t)

	newServer := NewServer( /* address */ testIp /* port */, testPort)
	newServer.ReadRow(&bytesOfTypicalResponse)

	ValidateFields(t, newServer)
}

func TestUpdateOnce(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUdpConn := NewMockUDPconnection(mockCtrl)
	mockUdpConn.EXPECT().Write(gomock.AssignableToTypeOf([]byte("s"))).Return(1, nil)

	bytesOfTypicalResponse := GetTypicalBytes(t)
	readingUdp := mockUdpConn.EXPECT().ReadFromUDP(gomock.AssignableToTypeOf([]byte(""))).DoAndReturn(func(b []byte) (_ int, _ *net.UDPAddr, _ error) {
		copy(b, bytesOfTypicalResponse)
		return len(bytesOfTypicalResponse), nil, nil
	}).Times(1)

	mockUdpConn.EXPECT().Close().Times(1).After(readingUdp)

	newServer := NewServer( /* address */ testIp /* port */, testPort)
	newServer.connection = mockUdpConn

	err := newServer.UpdateOnce()
	assert.NoError(t, err)

	ValidateFields(t, newServer)
}

func TestGetJoinLink(t *testing.T) {
	link := fmt.Sprintf("mtasa://%s:%d", testIp, testPort)

	testServer := NewServer(testIp, testPort)
	assert.Equal(t, link, testServer.GetJoinLink(), "join link supposed to contain ip and port of game server")
}
