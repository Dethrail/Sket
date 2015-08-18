package server

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"log"
)

func (p PacketResultJoinGame) ID() PacketID {
	return ID_ResultGameJoin
}

type PacketResultJoinGame struct {
	GameID uint32
	Result byte
}

func (p PacketResultJoinGame) Write(w *bufio.Writer) {
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(p.ID()))
	buf.WriteByte(0) // len first byte
	buf.WriteByte(0) // len second byte

	binary.Write(buf, binary.LittleEndian, p.GameID)
	buf.WriteByte(p.Result)

	var len uint16 = (uint16)(buf.Len())
	buf.Bytes()[1] = (byte)(len >> 0) // len first byte
	buf.Bytes()[2] = (byte)(len >> 8) // len second byte

	w.Write(buf.Bytes())
	err := w.Flush()

	if err != nil {
		log.Print("Error on result login" + err.Error() + "\n")
	}
}
