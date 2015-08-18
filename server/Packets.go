package server

import (
	//"log"
	"bufio"
	//"bytes"
	//"encoding/binary"
	"log"
	//"math"
	//"strings"
)

type PacketID uint8

const (
	ID_None = iota
	ID_Login

	ID_ResolveUDP

	ID_RequestGames
	ID_CreateGame
	ID_JoinGame
	ID_LeveGame
	ID_AbortGame

	ID_GOState
	ID_InstanceGO

	ID_PlayerInput

	ID_CharacterState

	ID_ResultLogin
	ID_ResultGameJoin
	ID_ResultInstanceGO
	ID_ResultRequestGames
	ID_Count
	//ID_ = PacketID(iota)

	//ID_Welcome    = PacketID(iota)
	//ID_EnterGame  = PacketID(iota)
	//ID_LoginError = PacketID(iota)

	//ID_SpawnPlayer     = PacketID(iota)
	//ID_PlayerInfo      = PacketID(iota)
	//ID_PlayerTransform = PacketID(iota)
	//ID_RemovePlayer    = PacketID(iota)
	//ID_PlayerMove      = PacketID(iota)

	//ID_Respawn = PacketID(iota)
)

type Packet interface {
	ID() PacketID
	Write(w *bufio.Writer)
}

// interface
func (p PacketLogin) Write(w *bufio.Writer)        {}
func (p PacketPlayerInput) Write(w *bufio.Writer)  {}
func (p PacketCreateGame) Write(w *bufio.Writer)   {}
func (p PacketJoinGame) Write(w *bufio.Writer)     {}
func (p PacketResolveUPD) Write(w *bufio.Writer)   {}
func (p PacketInstanceGO) Write(w *bufio.Writer)   {}
func (p PacketRequestGames) Write(w *bufio.Writer) {}

// start [Login]
func (p PacketLogin) ID() PacketID {
	return ID_Login
}

type PacketLogin struct {
	Name     string
	Password string
	State    byte
}

func ProcessPacketLogin(data []byte) Packet {
	var (
		name     string = ""
		password string = ""
		state    byte   = 0
	)

	id := data[0]
	if id != ID_Login {
		log.Printf("id of  packet (%d) != ID_Login(%d)", id, ID_Login)
	}

	strLength := data[3]
	var dataLength byte = 4 // [0] - id; [1,2] - packet lenght; [3] - str length; [4, 4 + strLength] - string

	name = string(data[dataLength : dataLength+strLength])
	dataLength += strLength

	if (byte)(len(data)) > dataLength {
		strLength = data[dataLength]
		dataLength++

		password = string(data[dataLength : dataLength+strLength])
		dataLength += strLength
	}
	return PacketLogin{name, password, state}
}

type PacketResultLogin struct {
	ClientID uint32
	State    byte
}

func (p PacketResultLogin) ID() PacketID {
	return ID_ResultLogin
}

// end [LoginResult]

// start [GameCreate]
func (p PacketCreateGame) ID() PacketID {
	return ID_CreateGame
}

type PacketCreateGame struct {
	HostPlayerID int32
	GameName     string
	IsPrivate    byte
	//Port         uint64
	State byte
}

func ProcessPacketCreateGame(data []byte) Packet {
	var (
		GameName  string = ""
		IsPrivate byte   = 0
		State     byte   = 0
	)

	id := data[0]
	if id != ID_CreateGame {
		log.Printf("id of  packet (%d) != ID_CreateGame(%d)", id, ID_CreateGame)
	}
	//packetLength = uint16(data[1]) | (uint16(data[2]) << 8)
	strLength := data[3]
	var dataLength byte = 4 // [0] - id; [1,2] - packet lenght; [3] - str length; [4, 4 + strLength] - string

	GameName = string(data[dataLength : dataLength+strLength])
	dataLength += strLength

	log.Printf("GameName = %s", GameName)
	IsPrivate = data[strLength+1]

	return PacketCreateGame{0, GameName, IsPrivate, State}
}

// end [GameCreate]

//type Welcome struct {
//	Name string
//}

//func (w Welcome) ID() PacketID {
//	return ID_Welcome
//}

//func NewWelcome(name string) Packet {
//	return Welcome{name}
//}

//type EnterGame struct {
//	PlayerID ID
//	Name     string
//}

//func (e EnterGame) ID() PacketID {
//	return ID_EnterGame
//}

//func NewEnterGame(id ID, name string) Packet {
//	return EnterGame{id, name}
//}

//type LoginError struct {
//	Error string
//}

//func (e LoginError) ID() PacketID {
//	return ID_LoginError
//}

//func NewLoginError(error string) Packet {
//	return LoginError{error}
//}

//type SpawnPlayer struct {
//	PlayerTransform
//	PlayerInfo
//}

//func (s SpawnPlayer) ID() PacketID {
//	return ID_SpawnPlayer
//}

//func NewSpawnPlayer(playerTransform PlayerTransform, playerInfo PlayerInfo) Packet {
//	return SpawnPlayer{playerTransform, playerInfo}
//}

//type PlayerTransform struct {
//	PlayerID ID
//	X, Y     float32
//	Rotation float32
//}

//func (s PlayerTransform) ID() PacketID {
//	return ID_PlayerTransform
//}

//func NewPlayerTransform(playerID ID, X, Y, Rotation float32) PlayerTransform {
//	return PlayerTransform{playerID, X, Y, Rotation}
//}

//type PlayerInfo struct {
//	PlayerID ID
//	Name     string
//}

//func (s PlayerInfo) ID() PacketID {
//	return ID_PlayerInfo
//}

//func NewPlayerInfo(playerID ID, name string) PlayerInfo {
//	return PlayerInfo{playerID, name}
//}

//type RemovePlayer struct {
//	PlayerID ID
//}

//func (s RemovePlayer) ID() PacketID {
//	return ID_RemovePlayer
//}

//func NewRemovePlayer(playerID ID) Packet {
//	return RemovePlayer{playerID}
//}

//type PlayerMove struct {
//	PlayerTransform
//}

//func (s PlayerMove) ID() PacketID {
//	return ID_PlayerMove
//}

//func NewPlayerMove(transfrom PlayerTransform) Packet {
//	return PlayerMove{transfrom}
//}

//type PlayerRespawn struct {
//}

//func (s PlayerRespawn) ID() PacketID {
//	return ID_Respawn
//}

//func NewPlayerRespawn() Packet {
//	return PlayerRespawn{}
//}
