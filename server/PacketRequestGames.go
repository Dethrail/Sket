package server

import (
//"log"
)

func (p PacketRequestGames) ID() PacketID {
	return ID_RequestGames
}

type PacketRequestGames struct{}

func ProcessPacketRequestGames(data []byte) Packet {
	return PacketRequestGames{}
}

func HandlePacketRequestGames(c *Client, rawPacket Packet) {
	//var p PacketRequestGames = rawPacket.(PacketRequestGames)
	var count = len(Games)
	var GamesID []uint32 = make([]uint32, count)
	for i, g := range Games {
		GamesID[i] = g.ID
	}
	answer := PacketResultRequestGames{GamesID, 0}
	answer.Write(c.TCPWriter)
	//if Games[p.GameId] != nil {
	//	Games[p.GameId].RegisterPlayer(c)
	//	log.Printf("Join to game (game not nil + %d)", len(Games[p.GameId].Players))
	//	answer := PacketResultGameJoin{p.GameId, 0}
	//	answer.Write(c.TCPWriter)
	//	for e := Games[p.GameId].Packages.Front(); e != nil; e = e.Next() {
	//		w := c.TCPWriter
	//		//log.Printf("owner[%d]", e.Value.(PacketResultInstanceGO).Owner)
	//		switch e.Value.(Packet).ID() {
	//		case ID_ResultInstanceGO:
	//			{
	//				var (
	//					GOID     uint32 = e.Value.(PacketResultInstanceGO).GOID
	//					Owner    uint32 = e.Value.(PacketResultInstanceGO).Owner
	//					x        float32
	//					y        float32
	//					result   byte = e.Value.(PacketResultInstanceGO).Result
	//					GOInGame bool = false
	//				)

	//				for _, obj := range Games[p.GameId].GameObjects {
	//					if obj.GOID == GOID {
	//						x = obj.Transform().Position().X
	//						y = obj.Transform().Position().Y
	//						GOInGame = true
	//					}
	//				}
	//				if GOInGame {
	//					//log.Printf("Send goid=%d and owner=%d", GOID, Owner)
	//					answer := PacketResultInstanceGO{GOID, Owner, x, y, result}
	//					answer.Write(w)
	//				}
	//			}
	//		default:
	//			{
	//				e.Value.(Packet).Write(w)
	//			}
	//		}
	//	}
	//} else {
	//	answer := PacketResultGameJoin{p.GameId, 1}
	//	answer.Write(c.TCPWriter)
	//}
}
