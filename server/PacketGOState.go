package server

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	//"runtime"
)

func (p PacketGOState) ID() PacketID {
	return ID_GOState
}

type PacketGOState struct {
	GOID      uint32
	X         float32
	Y         float32
	VX        float32
	VY        float32
	RotY      uint16
	Timestamp int32
}

func (p PacketGOState) Write(conn *net.UDPConn, addr *net.UDPAddr) {
	buf := new(bytes.Buffer)

	buf.WriteByte(byte(p.ID()))
	binary.Write(buf, binary.LittleEndian, p.GOID)
	binary.Write(buf, binary.LittleEndian, p.X)
	binary.Write(buf, binary.LittleEndian, p.Y)
	binary.Write(buf, binary.LittleEndian, p.VX)
	binary.Write(buf, binary.LittleEndian, p.VY)
	binary.Write(buf, binary.LittleEndian, p.RotY)
	binary.Write(buf, binary.LittleEndian, p.Timestamp)

	if conn == nil {
		log.Printf("udp connect is nil %t", (conn == nil))
		return
	}
	_, err := conn.WriteToUDP(buf.Bytes(), addr)

	//buf.Reset()
	//runtime.GC()
	if err != nil {
		log.Print("Error on PacketGOState" + err.Error() + "\n")
	}
}
