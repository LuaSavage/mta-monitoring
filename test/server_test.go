package server_test

import (
	"mta-monitoring/server"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadRow(t *testing.T) {

	// responce from ase udp
	typicalResponse := `"EYE1\x04mta\x0622044M                          MTA:SA Türkiye - Norm Gaming [ Turkish / Turkey ]\aMTA:SA\x05None\x041.5\x021\x020\x0370\x01"`
	testServer := &server.Server{Address: "217.106.106.107", Port: 22044, AsePort: 22044 + 123}
	unquotedTypicalResponse, _ := strconv.Unquote(typicalResponse)
	bytesOfTypicalResponse := []byte(unquotedTypicalResponse)

	testServer.ReadRow(&bytesOfTypicalResponse)

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
		assert.Equal(t, field.Interface(), value, "Field "+key+" should be equal")
	}
}

func TestGetJoinLink(t *testing.T) {
	address := "217.106.106.107"
	port := 22044
	testServer := &server.Server{Address: address, Port: port, AsePort: port + 123}
	assert.Equal(t, `mtasa://`+address+`:`+strconv.Itoa(port), testServer.GetJoinLink(), "join link supposed to contain ip and port of game server")
}
