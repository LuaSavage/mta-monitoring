package server

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const typicalResponse = `"EYE1\x04mta\x0622044M                          MTA:SA Türkiye - Norm Gaming [ Turkish / Turkey ]\aMTA:SA\x05None\x041.5\x021\x020\x0370\x01"`
const testIp = "217.106.106.107"
const testPort = 22044

func GetTypicalBytes(t *testing.T) []byte {
	t.Helper()

	unquotedTypicalResponse, _ := strconv.Unquote(typicalResponse)
	return []byte(unquotedTypicalResponse)
}

func ValidateFields(t *testing.T, testServer *Server) {
	t.Helper()

	expectedValues := map[string]interface{}{
		"Game":       "mta",
		"Port":       22044,
		"Name":       "                          MTA:SA Türkiye - Norm Gaming [ Turkish / Turkey ]",
		"Gamemode":   "MTA:SA",
		"Map":        "None",
		"Version":    "1.5",
		"Somewhat":   "1",
		"Players":    0,
		"Maxplayers": 70,
	}

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
	mockUdpConn.EXPECT().ReadFromUDP(gomock.AssignableToTypeOf(emptyByte)).DoAndReturn(func(b []byte) (int, *net.UDPAddr, error) {
		copy(b, bytesOfTypicalResponse)
		return len(bytesOfTypicalResponse), nil, nil
	}).Times(1)

	newServer := NewServer(testIp, testPort)
	newServer.connection = mockUdpConn
	receivedBytes, err := newServer.ReadSocketData()

	assert.NoError(t, err)
	receivedBytesShort := (*receivedBytes)[:len(bytesOfTypicalResponse)]
	assert.True(t, bytes.Equal(bytesOfTypicalResponse, receivedBytesShort), "Received bytes unequal")
}

func TestReadRow(t *testing.T) {
	bytesOfTypicalResponse := GetTypicalBytes(t)

	newServer := NewServer(testIp, testPort)
	newServer.ReadRow(&bytesOfTypicalResponse)

	ValidateFields(t, newServer)
}

func TestUpdateOnce(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUdpConn := NewMockUDPconnection(mockCtrl)
	mockUdpConn.EXPECT().Write(gomock.AssignableToTypeOf([]byte("s"))).Return(1, nil)

	bytesOfTypicalResponse := GetTypicalBytes(t)
	readingUdp := mockUdpConn.EXPECT().ReadFromUDP(gomock.AssignableToTypeOf([]byte(""))).DoAndReturn(func(b []byte) (int, *net.UDPAddr, error) {
		copy(b, bytesOfTypicalResponse)
		return len(bytesOfTypicalResponse), nil, nil
	}).Times(1)

	mockUdpConn.EXPECT().Close().Times(1).After(readingUdp)

	newServer := NewServer(testIp, testPort)
	newServer.connection = mockUdpConn

	err := newServer.UpdateOnce()
	assert.NoError(t, err)

	ValidateFields(t, newServer)
}

func TestUpdateOnce_ErrorOnWrite(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUdpConn := NewMockUDPconnection(mockCtrl)

	connWriteErr := errors.New("smth wrong with connection write")
	mockUdpConn.EXPECT().Write(gomock.AssignableToTypeOf([]byte("s"))).Return(1, connWriteErr)
	mockUdpConn.EXPECT().Close().Times(1)

	newServer := NewServer(testIp, testPort)
	newServer.connection = mockUdpConn

	err := newServer.UpdateOnce()
	assert.Error(t, err)
}

func TestUpdateOnce_ErrOnReadingUDP(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUdpConn := NewMockUDPconnection(mockCtrl)

	udpReadingErr := errors.New("smth wrong with udp reading")
	mockUdpConn.EXPECT().Write(gomock.AssignableToTypeOf([]byte("s"))).Return(1, nil)

	bytesOfTypicalResponse := GetTypicalBytes(t)
	readingUdp := mockUdpConn.EXPECT().ReadFromUDP(gomock.AssignableToTypeOf([]byte(""))).DoAndReturn(func(b []byte) (int, *net.UDPAddr, error) {
		copy(b, bytesOfTypicalResponse)
		return len(bytesOfTypicalResponse), nil, udpReadingErr
	}).Times(1)

	mockUdpConn.EXPECT().Close().Times(1).After(readingUdp)

	newServer := NewServer(testIp, testPort)
	newServer.connection = mockUdpConn

	err := newServer.UpdateOnce()
	assert.Error(t, err)
}

func TestUpdateOnce_ReadTimeout(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUdpConn := NewMockUDPconnection(mockCtrl)

	timeoutErr := &net.OpError{
		Op:  "read",
		Net: "udp",
		Err: os.ErrDeadlineExceeded,
	}

	mockUdpConn.EXPECT().Write(gomock.AssignableToTypeOf([]byte("s"))).Return(1, nil)
	readingUdp := mockUdpConn.EXPECT().ReadFromUDP(gomock.AssignableToTypeOf([]byte(""))).Return(0, nil, timeoutErr).Times(1)
	mockUdpConn.EXPECT().Close().Times(1).After(readingUdp)

	newServer := NewServer(testIp, testPort)
	newServer.connection = mockUdpConn

	err := newServer.UpdateOnce()
	assert.Error(t, err)
	assert.ErrorIs(t, err, os.ErrDeadlineExceeded)
}

func TestGetJoinLink(t *testing.T) {
	link := fmt.Sprintf("mtasa://%s:%d", testIp, testPort)

	testServer := NewServer(testIp, testPort)
	assert.Equal(t, link, testServer.GetJoinLink(), "Join link supposed to contain ip and port of game server")
}
