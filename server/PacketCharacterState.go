package server

import (
	"bufio"
	//"bytes"
	"encoding/binary"
	"log"
	//"math"
)

const (
	CS_Default = iota
	CS_Attack
	CS_Death
	CS_Revive
	CS_Destroy
)

func (p PacketCharacterState) ID() PacketID {
	return ID_CharacterState
}

type PacketCharacterState struct {
	GOID      uint32
	CharState byte
}

func (p PacketCharacterState) Write(w *bufio.Writer) {

	w.WriteByte(byte(p.ID()))
	binary.Write(w, binary.LittleEndian, p.GOID)
	w.WriteByte(p.CharState)
	//binary.Write(w, binary.LittleEndian, EndPacket)

	err := w.Flush()
	if err != nil {
		log.Print("Error on PacketCharacterState:" + err.Error() + "\n")
	}
}
