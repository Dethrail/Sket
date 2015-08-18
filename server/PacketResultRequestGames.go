package server

import (
	"bufio"
	"encoding/binary"
	"log"
)

func (p PacketResultRequestGames) ID() PacketID {
	return ID_ResultRequestGames
}

type PacketResultRequestGames struct {
	GamesID []uint32

	Result byte
}

func (p PacketResultRequestGames) Write(w *bufio.Writer) {
	w.WriteByte(byte(p.ID()))         // 1
	w.WriteByte(byte(len(p.GamesID))) // 2 len
	//log.Printf("len %d", byte(len(p.GamesID)))
	for _, id := range p.GamesID {
		binary.Write(w, binary.LittleEndian, id) // 2 + 4 * len
		//log.Printf("id = %d", id)
	}

	w.WriteByte(p.Result) // 1
	//binary.Write(w, binary.LittleEndian, EndPacket)

	err := w.Flush()
	if err != nil {
		log.Print("Error on PacketResultRequestGames" + err.Error() + "\n")
	}
}
