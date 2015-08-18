package server

import (
	"log"
)

func (p PacketPlayerInput) ID() PacketID {
	return ID_PlayerInput
}

type PacketPlayerInput struct {
	Key      byte
	KeyState byte
}

func ProcessPacketPlayerInput(data []byte) Packet {
	var (
		key      byte
		keyState byte
	)

	id := data[0]
	if id != ID_PlayerInput {
		log.Printf("id of packet (%d) != ID_PlayerInput(%d)", id, ID_PlayerInput)
	}
	//intVal := uint32(val[0]) | (uint32(val[1]) << 8) | (uint32(val[2])<<16 | uint32(val[3])<<24)

	//packetLength := uint16(data[1]) | (uint16(data[2]) << 8)
	key = data[3]
	keyState = data[4]

	return PacketPlayerInput{key, keyState}
}

func HandlePacketPlayerInput(c *Client, rawPacket Packet) {
	var p PacketPlayerInput = rawPacket.(PacketPlayerInput)

	if p.Key == 'a' {
		c.Character.Input.A = p.KeyState
		if p.KeyState == 1 {
			c.Character.Direction.X--
		} else {
			c.Character.Direction.X++
		}
	}

	if p.Key == 'd' {
		c.Character.Input.D = p.KeyState
		if p.KeyState == 1 {
			c.Character.Direction.X++
		} else {
			c.Character.Direction.X--
		}
	}

	if p.Key == 'w' {
		c.Character.Input.W = p.KeyState
		if p.KeyState == 1 {
			c.Character.Direction.Y++
		} else {
			c.Character.Direction.Y--
		}
	}

	if p.Key == 's' {
		c.Character.Input.D = p.KeyState
		if p.KeyState == 1 {
			c.Character.Direction.Y--
		} else {
			c.Character.Direction.Y++
		}
	}

	if p.Key == ' ' { // Space
		c.Character.Input.Space = p.KeyState
	}

	if p.Key == 1 { // shift
		c.Character.Input.Shift = p.KeyState
	}
	//log.Printf("input x=%d", p.Key)
}

// **** PACKET can only receive ****

//func (p PacketPlayerInput) Write(w *bufio.Writer) {
//}
