package server

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const typicalResponse = `"EYE1\x04mta\x0622044M                          MTA:SA Türkiye - Norm Gaming [ Turkish / Turkey ]\aMTA:SA\x05None\x041.5\x021\x020\x0370\x01"`

func typicalResponseBytes(t *testing.T) []byte {
	t.Helper()

	unquoted, err := strconv.Unquote(typicalResponse)
	require.NoError(t, err)
	return []byte(unquoted)
}

func writeLengthPrefixedString(buf *bytes.Buffer, value string) {
	buf.WriteByte(byte(len(value) + 1))
	buf.WriteString(value)
}

func buildResponseWithPlayers(t *testing.T) []byte {
	t.Helper()

	var buf bytes.Buffer
	buf.WriteString("EYE1")
	writeLengthPrefixedString(&buf, "mta")
	writeLengthPrefixedString(&buf, "22044")
	writeLengthPrefixedString(&buf, "Test Server")
	writeLengthPrefixedString(&buf, "MTA:SA")
	writeLengthPrefixedString(&buf, "None")
	writeLengthPrefixedString(&buf, "1.5")
	writeLengthPrefixedString(&buf, "0")
	writeLengthPrefixedString(&buf, "2")
	writeLengthPrefixedString(&buf, "70")
	buf.WriteByte(1)

	for _, player := range []Player{
		{Name: "Alice", Score: 10, Ping: 42},
		{Name: "Bob", Score: 5, Ping: 88},
	} {
		buf.WriteByte(0x3F)
		writeLengthPrefixedString(&buf, player.Name)
		buf.WriteByte(1)
		buf.WriteByte(1)
		writeLengthPrefixedString(&buf, strconv.Itoa(player.Score))
		writeLengthPrefixedString(&buf, strconv.Itoa(player.Ping))
		buf.WriteByte(1)
	}

	return buf.Bytes()
}

func buildInvalidPasswordedResponse(t *testing.T) []byte {
	t.Helper()

	var buf bytes.Buffer
	buf.WriteString("EYE1")
	writeLengthPrefixedString(&buf, "mta")
	writeLengthPrefixedString(&buf, "22044")
	writeLengthPrefixedString(&buf, "Test Server")
	writeLengthPrefixedString(&buf, "MTA:SA")
	writeLengthPrefixedString(&buf, "None")
	writeLengthPrefixedString(&buf, "1.5")
	writeLengthPrefixedString(&buf, "yes")
	return buf.Bytes()
}

func TestParseASE_typicalResponse(t *testing.T) {
	parsed, err := parseASE(typicalResponseBytes(t))
	require.NoError(t, err)

	assert.Equal(t, "mta", parsed.Game)
	assert.Equal(t, 22044, parsed.Port)
	assert.Equal(t, "                          MTA:SA Türkiye - Norm Gaming [ Turkish / Turkey ]", parsed.Name)
	assert.Equal(t, "MTA:SA", parsed.Gamemode)
	assert.Equal(t, "None", parsed.Map)
	assert.Equal(t, "1.5", parsed.Version)
	assert.True(t, parsed.Passworded)
	assert.Equal(t, 0, parsed.Players)
	assert.Equal(t, 70, parsed.Maxplayers)
	assert.Empty(t, parsed.PlayerList)
}

func TestParseASE_withPlayers(t *testing.T) {
	parsed, err := parseASE(buildResponseWithPlayers(t))
	require.NoError(t, err)

	assert.Equal(t, 2, parsed.Players)
	assert.Equal(t, []Player{
		{Name: "Alice", Score: 10, Ping: 42},
		{Name: "Bob", Score: 5, Ping: 88},
	}, parsed.PlayerList)
}

func TestParseASE_errors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    []byte
		wantErr string
	}{
		{
			name:    "invalid header",
			data:    []byte("BAD1\x04mta"),
			wantErr: "invalid header",
		},
		{
			name:    "truncated buffer",
			data:    []byte("EYE1\x04mt"),
			wantErr: "unexpected EOF",
		},
		{
			name:    "invalid passworded",
			data:    buildInvalidPasswordedResponse(t),
			wantErr: "invalid passworded value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := parseASE(tt.data)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}
