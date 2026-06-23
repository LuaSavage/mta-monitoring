package server

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

const aseHeader = "EYE1"

const (
	flagPlayerName  = 0x01
	flagPlayerTeam  = 0x02
	flagPlayerSkin  = 0x04
	flagPlayerScore = 0x08
	flagPlayerPing  = 0x10
	flagPlayerTime  = 0x20
)

type parsedASE struct {
	Game       string
	Port       int
	Name       string
	Gamemode   string
	Map        string
	Version    string
	Passworded bool
	Players    int
	Maxplayers int
	PlayerList []Player
}

func parseASE(data []byte) (parsedASE, error) {
	if len(data) < len(aseHeader) {
		return parsedASE{}, errors.New("unexpected EOF")
	}

	if string(data[:len(aseHeader)]) != aseHeader {
		return parsedASE{}, fmt.Errorf("invalid header: want %q", aseHeader)
	}

	buf := bytes.NewBuffer(data[len(aseHeader):])

	game, err := readLengthPrefixedString(buf)
	if err != nil {
		return parsedASE{}, fmt.Errorf("read game: %w", err)
	}

	portStr, err := readLengthPrefixedString(buf)
	if err != nil {
		return parsedASE{}, fmt.Errorf("read port: %w", err)
	}
	port, err := parseIntField(portStr, "port")
	if err != nil {
		return parsedASE{}, err
	}

	name, err := readLengthPrefixedString(buf)
	if err != nil {
		return parsedASE{}, fmt.Errorf("read name: %w", err)
	}

	gamemode, err := readLengthPrefixedString(buf)
	if err != nil {
		return parsedASE{}, fmt.Errorf("read gamemode: %w", err)
	}

	mapName, err := readLengthPrefixedString(buf)
	if err != nil {
		return parsedASE{}, fmt.Errorf("read map: %w", err)
	}

	version, err := readLengthPrefixedString(buf)
	if err != nil {
		return parsedASE{}, fmt.Errorf("read version: %w", err)
	}

	passwordedStr, err := readLengthPrefixedString(buf)
	if err != nil {
		return parsedASE{}, fmt.Errorf("read passworded: %w", err)
	}
	passworded, err := parsePassworded(passwordedStr)
	if err != nil {
		return parsedASE{}, err
	}

	playersStr, err := readLengthPrefixedString(buf)
	if err != nil {
		return parsedASE{}, fmt.Errorf("read players: %w", err)
	}
	players, err := parseIntField(playersStr, "players")
	if err != nil {
		return parsedASE{}, err
	}

	maxPlayersStr, err := readLengthPrefixedString(buf)
	if err != nil {
		return parsedASE{}, fmt.Errorf("read maxplayers: %w", err)
	}
	maxPlayers, err := parseIntField(maxPlayersStr, "maxplayers")
	if err != nil {
		return parsedASE{}, err
	}

	if err := skipRules(buf); err != nil {
		return parsedASE{}, fmt.Errorf("read rules: %w", err)
	}

	playerList, err := parsePlayers(buf)
	if err != nil {
		return parsedASE{}, fmt.Errorf("read player list: %w", err)
	}

	return parsedASE{
		Game:       game,
		Port:       port,
		Name:       name,
		Gamemode:   gamemode,
		Map:        mapName,
		Version:    version,
		Passworded: passworded,
		Players:    players,
		Maxplayers: maxPlayers,
		PlayerList: playerList,
	}, nil
}

func readLengthPrefixedString(buf *bytes.Buffer) (string, error) {
	lengthByte, err := buf.ReadByte()
	if err != nil {
		return "", fmt.Errorf("read length byte: %w", err)
	}

	payloadLen := int(lengthByte) - 1
	if payloadLen < 0 {
		return "", errors.New("invalid string length")
	}
	if buf.Len() < payloadLen {
		return "", errors.New("unexpected EOF")
	}

	return string(buf.Next(payloadLen)), nil
}

func skipLengthPrefixedString(buf *bytes.Buffer) error {
	_, err := readLengthPrefixedString(buf)
	return err
}

func parseIntField(value string, field string) (int, error) {
	n, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", field, err)
	}
	return n, nil
}

func parsePassworded(value string) (bool, error) {
	switch value {
	case "0":
		return false, nil
	case "1":
		return true, nil
	default:
		return false, fmt.Errorf("invalid passworded value: %q", value)
	}
}

func skipRules(buf *bytes.Buffer) error {
	for {
		keyLenByte, err := buf.ReadByte()
		if err != nil {
			return fmt.Errorf("read rule key length: %w", err)
		}
		if keyLenByte == 1 {
			return nil
		}

		keyLen := int(keyLenByte) - 1
		if buf.Len() < keyLen {
			return errors.New("unexpected EOF")
		}
		buf.Next(keyLen)

		if err := skipLengthPrefixedString(buf); err != nil {
			return fmt.Errorf("read rule value: %w", err)
		}
	}
}

func parsePlayers(buf *bytes.Buffer) ([]Player, error) {
	var players []Player

	for buf.Len() > 0 {
		flags, err := buf.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("read player flags: %w", err)
		}

		player, err := parsePlayer(buf, flags)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}

	return players, nil
}

func parsePlayer(buf *bytes.Buffer, flags byte) (Player, error) {
	var player Player

	if flags&flagPlayerName != 0 {
		name, err := readLengthPrefixedString(buf)
		if err != nil {
			return Player{}, fmt.Errorf("read player name: %w", err)
		}
		player.Name = name
	}

	if flags&flagPlayerTeam != 0 {
		if err := skipLengthPrefixedString(buf); err != nil {
			return Player{}, fmt.Errorf("skip player team: %w", err)
		}
	}

	if flags&flagPlayerSkin != 0 {
		if err := skipLengthPrefixedString(buf); err != nil {
			return Player{}, fmt.Errorf("skip player skin: %w", err)
		}
	}

	if flags&flagPlayerScore != 0 {
		scoreStr, err := readLengthPrefixedString(buf)
		if err != nil {
			return Player{}, fmt.Errorf("read player score: %w", err)
		}
		score, err := parseIntField(scoreStr, "player score")
		if err != nil {
			return Player{}, err
		}
		player.Score = score
	}

	if flags&flagPlayerPing != 0 {
		pingStr, err := readLengthPrefixedString(buf)
		if err != nil {
			return Player{}, fmt.Errorf("read player ping: %w", err)
		}
		ping, err := parseIntField(pingStr, "player ping")
		if err != nil {
			return Player{}, err
		}
		player.Ping = ping
	}

	if flags&flagPlayerTime != 0 {
		if err := skipLengthPrefixedString(buf); err != nil {
			return Player{}, fmt.Errorf("skip player time: %w", err)
		}
	}

	return player, nil
}

func applyParsedASE(s *Server, parsed parsedASE) {
	s.Game = parsed.Game
	s.Port = parsed.Port
	s.Name = parsed.Name
	s.Gamemode = parsed.Gamemode
	s.Map = parsed.Map
	s.Version = parsed.Version
	s.Passworded = parsed.Passworded
	s.Players = parsed.Players
	s.Maxplayers = parsed.Maxplayers
	s.PlayerList = parsed.PlayerList
}
