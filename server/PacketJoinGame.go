package server

import (
	"log"
)

func (p PacketJoinGame) ID() PacketID {
	return ID_JoinGame
}

type PacketJoinGame struct {
	GameId   uint32
	GameName string
}

func ProcessPacketJoinGame(data []byte) Packet {
	var (
		gameId   uint32 = 0
		gameName string = ""
	)

	//rawData := bytes.Split(data, []byte{129})
	//for i, val := range rawData {
	//	switch i {
	//	case 0:
	//		{
	id := data[0]
	if id != ID_JoinGame {
		log.Printf("id of  packet (%d) != ID_JoinGame(%d)", id, ID_JoinGame)
	}
	//}
	//case 1:
	//{
	//packetLength := uint16(data[1]) | (uint16(data[2]) << 8)
	gameId = uint32(data[3]) | (uint32(data[4]) << 8) | (uint32(data[5])<<16 | uint32(data[6])<<24)
	//if game id = 0, client don't know game id
	//log.Printf("game id(%d)", gameId)
	//}
	//case 2:
	//{
	var dataLenght byte = 7
	strLenght := data[dataLenght]
	dataLenght++

	gameName = string(data[dataLenght : dataLenght+strLenght])
	//log.Printf("game name(%s)", gameName)
	//}

	//default:
	//{
	//	log.Printf("parameter #%d = first byte=%d", i, val[0])
	//}
	//}
	//}

	//for i, val := range arr {

	return PacketJoinGame{gameId, gameName}
}

func HandlePacketJoinGame(c *Client, rawPacket Packet) {
	var p PacketJoinGame = rawPacket.(PacketJoinGame)
	if Games[p.GameId] != nil {
		Games[p.GameId].RegisterPlayer(c)
		log.Printf("Join to game (game not nil + %d)", len(Games[p.GameId].Players))
		answer := PacketResultJoinGame{p.GameId, 0}
		answer.Write(c.TCPWriter)
		for e := Games[p.GameId].Packages.Front(); e != nil; e = e.Next() {
			w := c.TCPWriter
			//log.Printf("owner[%d]", e.Value.(PacketResultInstanceGO).Owner)
			switch e.Value.(Packet).ID() {
			case ID_ResultInstanceGO:
				{
					var (
						GOID     uint32 = e.Value.(PacketResultInstanceGO).GOID
						Owner    uint32 = e.Value.(PacketResultInstanceGO).Owner
						x        float32
						y        float32
						result   byte = e.Value.(PacketResultInstanceGO).Result
						GOInGame bool = false
					)

					for _, obj := range Games[p.GameId].GameObjects {
						if obj.GOID == GOID {
							x = obj.Transform().Position().X
							y = obj.Transform().Position().Y
							GOInGame = true
						}
					}
					if GOInGame {
						//log.Printf("Send goid=%d and owner=%d", GOID, Owner)
						answer := PacketResultInstanceGO{GOID, Owner, x, y, result}
						answer.Write(w)
					}
				}
			default:
				{
					e.Value.(Packet).Write(w)
				}
			}
		}
	} else {
		answer := PacketResultJoinGame{p.GameId, 1}
		answer.Write(c.TCPWriter)
	}
}
