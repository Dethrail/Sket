package server

import (
	"bufio"
	"encoding/binary"
	"log"
)

func (p PacketResultInstanceGO) ID() PacketID {
	return ID_ResultInstanceGO
}

type PacketResultInstanceGO struct {
	GOID   uint32
	Owner  uint32
	X      float32
	Y      float32
	Result byte
}

func (p PacketResultInstanceGO) Write(w *bufio.Writer) {
	w.WriteByte(byte(p.ID()))                     // 1
	binary.Write(w, binary.LittleEndian, p.GOID)  // 4
	binary.Write(w, binary.LittleEndian, p.Owner) // 4
	binary.Write(w, binary.LittleEndian, p.X)     // 4
	binary.Write(w, binary.LittleEndian, p.Y)     // 4
	w.WriteByte(p.Result)                         // 1
	//binary.Write(w, binary.LittleEndian, EndPacket)

	err := w.Flush()
	if err != nil {
		log.Print("Error on PacketResultInstanceGO" + err.Error() + "\n")
	}
}
