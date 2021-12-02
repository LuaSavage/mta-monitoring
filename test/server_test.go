package server_test

import (
	"mta-monitoring/server"
	"testing"
  "fmt"
	"github.com/stretchr/testify/assert"
)

func TestReadRow(t *testing.T) {

  // responce from ase udp

  response := []byte(`EYE1\x04mta\x0622003][BS][RU] MTA DayZ Ultimate [#1] "Brown" [HARDCORE, LOOT X1, 130 VEHICLES, FAST ZOMBIE] -TOP-\x14DayZ Ultimate 1.4.7\x05None\x041.5\x020\x024\x0340\x01?\x10[USSR]BadBoy228\x01\x01\x01\x0370\x01?\x07Lesnik\x01\x01\x01\x0332\x01?\tVendetta\x01\x01\x01\x0379\x01?\x05Netx\x01\x01\x01\x0395\x01`) 
  
    fmt.Println(response)
  //newServer := server{}
  newServer := server.NewServer("217.106.106.107", 22044)

  assert.Equal(t, newServer.Port, 22044)


  newServer.ReadRow(&response)
/*
  // assert equality
  assert.Equal(t, newServer.Game, "mta")
  assert.Equal(t, newServer.Port, 22003)
  assert.Equal(t, newServer.Name, `[BS][RU] MTA DayZ Ultimate [#1] "Brown" '
                                           '[HARDCORE, LOOT X1, 130 VEHICLES, FAST ZOMBIE] -TOP-`)

  assert.Equal(t, newServer.Gamemode, "DayZ Ultimate 1.4.7")
  assert.Equal(t, newServer.Map, "None")
  assert.Equal(t, newServer.Version, "1.5")
  assert.Equal(t, newServer.Somewhat, '0')
  assert.Equal(t, newServer.Players, 4)
  assert.Equal(t, newServer.Maxplayers, 40)

*/
  // assert inequality
  //assert.NotEqual(t, 123, 456, "they should not be equal")

}


