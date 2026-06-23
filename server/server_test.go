package server

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const testIP = "217.106.106.107"
const testPort = 22044

func validateTypicalFields(t *testing.T, testServer *Server) {
	t.Helper()

	assert.Equal(t, "mta", testServer.Game)
	assert.Equal(t, 22044, testServer.Port)
	assert.Equal(t, "                          MTA:SA Türkiye - Norm Gaming [ Turkish / Turkey ]", testServer.Name)
	assert.Equal(t, "MTA:SA", testServer.Gamemode)
	assert.Equal(t, "None", testServer.Map)
	assert.Equal(t, "1.5", testServer.Version)
	assert.True(t, testServer.Passworded)
	assert.Equal(t, 0, testServer.Players)
	assert.Equal(t, 70, testServer.Maxplayers)
	assert.Empty(t, testServer.PlayerList)
}

func TestReadSocketData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUDPConn := NewMockUDPconnection(mockCtrl)
	mockUDPConn.EXPECT().Write(gomock.AssignableToTypeOf([]byte("s"))).Return(1, nil)

	emptyByte := []byte("")

	bytesOfTypicalResponse := typicalResponseBytes(t)
	mockUDPConn.EXPECT().ReadFromUDP(gomock.AssignableToTypeOf(emptyByte)).DoAndReturn(func(b []byte) (int, *net.UDPAddr, error) {
		copy(b, bytesOfTypicalResponse)
		return len(bytesOfTypicalResponse), nil, nil
	}).Times(1)

	newServer := NewServer(testIP, testPort)
	newServer.connection = mockUDPConn
	receivedBytes, err := newServer.ReadSocketData()

	assert.NoError(t, err)
	receivedBytesShort := (*receivedBytes)[:len(bytesOfTypicalResponse)]
	assert.True(t, bytes.Equal(bytesOfTypicalResponse, receivedBytesShort), "Received bytes unequal")
}

func TestReadRow(t *testing.T) {
	bytesOfTypicalResponse := typicalResponseBytes(t)

	newServer := NewServer(testIP, testPort)
	err := newServer.ReadRow(&bytesOfTypicalResponse)
	require.NoError(t, err)

	validateTypicalFields(t, newServer)
}

func TestReadRow_withPlayers(t *testing.T) {
	response := buildResponseWithPlayers(t)

	newServer := NewServer(testIP, testPort)
	err := newServer.ReadRow(&response)
	require.NoError(t, err)

	assert.Equal(t, []Player{
		{Name: "Alice", Score: 10, Ping: 42},
		{Name: "Bob", Score: 5, Ping: 88},
	}, newServer.PlayerList)
}

func TestUpdateOnce(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUDPConn := NewMockUDPconnection(mockCtrl)
	mockUDPConn.EXPECT().Write(gomock.AssignableToTypeOf([]byte("s"))).Return(1, nil)

	bytesOfTypicalResponse := typicalResponseBytes(t)
	readingUDP := mockUDPConn.EXPECT().ReadFromUDP(gomock.AssignableToTypeOf([]byte(""))).DoAndReturn(func(b []byte) (int, *net.UDPAddr, error) {
		copy(b, bytesOfTypicalResponse)
		return len(bytesOfTypicalResponse), nil, nil
	}).Times(1)

	mockUDPConn.EXPECT().Close().Times(1).After(readingUDP)

	newServer := NewServer(testIP, testPort)
	newServer.connection = mockUDPConn

	err := newServer.UpdateOnce()
	assert.NoError(t, err)

	validateTypicalFields(t, newServer)
}

func TestUpdateOnce_ErrorOnWrite(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUDPConn := NewMockUDPconnection(mockCtrl)

	connWriteErr := errors.New("smth wrong with connection write")
	mockUDPConn.EXPECT().Write(gomock.AssignableToTypeOf([]byte("s"))).Return(1, connWriteErr)
	mockUDPConn.EXPECT().Close().Times(1)

	newServer := NewServer(testIP, testPort)
	newServer.connection = mockUDPConn

	err := newServer.UpdateOnce()
	assert.Error(t, err)
}

func TestUpdateOnce_ErrOnReadingUDP(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUDPConn := NewMockUDPconnection(mockCtrl)

	udpReadingErr := errors.New("smth wrong with udp reading")
	mockUDPConn.EXPECT().Write(gomock.AssignableToTypeOf([]byte("s"))).Return(1, nil)

	bytesOfTypicalResponse := typicalResponseBytes(t)
	readingUDP := mockUDPConn.EXPECT().ReadFromUDP(gomock.AssignableToTypeOf([]byte(""))).DoAndReturn(func(b []byte) (int, *net.UDPAddr, error) {
		copy(b, bytesOfTypicalResponse)
		return len(bytesOfTypicalResponse), nil, udpReadingErr
	}).Times(1)

	mockUDPConn.EXPECT().Close().Times(1).After(readingUDP)

	newServer := NewServer(testIP, testPort)
	newServer.connection = mockUDPConn

	err := newServer.UpdateOnce()
	assert.Error(t, err)
}

func TestUpdateOnce_ReadTimeout(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUDPConn := NewMockUDPconnection(mockCtrl)

	timeoutErr := &net.OpError{
		Op:  "read",
		Net: "udp",
		Err: os.ErrDeadlineExceeded,
	}

	mockUDPConn.EXPECT().Write(gomock.AssignableToTypeOf([]byte("s"))).Return(1, nil)
	readingUDP := mockUDPConn.EXPECT().ReadFromUDP(gomock.AssignableToTypeOf([]byte(""))).Return(0, nil, timeoutErr).Times(1)
	mockUDPConn.EXPECT().Close().Times(1).After(readingUDP)

	newServer := NewServer(testIP, testPort)
	newServer.connection = mockUDPConn

	err := newServer.UpdateOnce()
	assert.Error(t, err)
	assert.ErrorIs(t, err, os.ErrDeadlineExceeded)
}

func TestGetJoinLink(t *testing.T) {
	link := fmt.Sprintf("mtasa://%s:%d", testIP, testPort)

	testServer := NewServer(testIP, testPort)
	assert.Equal(t, link, testServer.GetJoinLink(), "Join link supposed to contain ip and port of game server")
}
