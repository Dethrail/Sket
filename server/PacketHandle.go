package server

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"github.com/Dethrail/chipmunk"
	"github.com/Dethrail/chipmunk/vect"
	"log"
)

func HandlePacketLogin(c *Client, rawPacket Packet) {
	var p PacketLogin = rawPacket.(PacketLogin)
	log.Printf("name = %s\n", p.Name)
	log.Printf("password = % s\n", p.Password)
	//if p.Name == "Max" && p.Password == "0263541zxc" {
	p.State = 0
	//} else {
	//	p.State = 1
	//}

	answer := PacketResultLogin{uint32(c.ID), p.State}
	answer.Write(c.TCPWriter)
}

func (p PacketResultLogin) Write(w *bufio.Writer) {
	buf := new(bytes.Buffer)

	buf.WriteByte(byte(p.ID()))
	buf.WriteByte(0) // len first byte
	buf.WriteByte(0) // len second byte

	binary.Write(buf, binary.LittleEndian, p.ClientID)
	buf.WriteByte(p.State) // logged in

	var len uint16 = (uint16)(buf.Len())
	buf.Bytes()[1] = (byte)(len >> 0) // len first byte
	buf.Bytes()[2] = (byte)(len >> 8) // len second byte

	//log.Printf("%d", buf.Bytes()[1])
	//log.Printf("%d", buf.Bytes()[2])

	w.Write(buf.Bytes())
	err := w.Flush()

	if err != nil {
		log.Print("Error on result login" + err.Error() + "\n")
	}
}

func HandlePacketCreateGame(c *Client, rawPacket Packet) { // c *Client,
	//var p PacketGameCreate = rawPacket.(PacketGameCreate)
	g, id := CreateGame(c)

	floor := NewGameObject("ground")
	floor.AddComponent(NewPhysicsShape(true, g.Space, chipmunk.NewBox(vect.Vector_Zero, 200, 0.3)))
	////floor.Physics.Body.SetMoment(Inf)
	//floor.Physics.Shape.SetFriction(0)
	//floor.Physics.Shape.SetElasticity(0)
	//floor.Transform().SetPositionf(0, -0.15)
	g.AddGameObject(floor)

	//wallLeft := NewGameObject("wall left")
	//wallLeft.AddComponent(NewPhysicsShape(true, g.Space, chipmunk.NewBox(vect.Vector_Zero, 1, 100)))
	////floor.Physics.Body.SetMoment(Inf)
	//wallLeft.Physics.Shape.SetFriction(0)
	//wallLeft.Physics.Shape.SetElasticity(0)
	//wallLeft.Transform().SetPositionf(-9.4, 0)
	//g.AddGameObject(wallLeft)

	//wallRight := NewGameObject("wall right")
	//wallRight.AddComponent(NewPhysicsShape(true, g.Space, chipmunk.NewBox(vect.Vector_Zero, 1, 100)))
	////floor.Physics.Body.SetMoment(Inf)
	//wallRight.Physics.Shape.SetFriction(0)
	//wallRight.Physics.Shape.SetElasticity(0)
	//wallRight.Transform().SetPositionf(9.4, 0)
	//g.AddGameObject(wallRight)

	//platform := NewGameObject("platform 1")
	//platform.AddComponent(NewPhysicsShape(true, g.Space, chipmunk.NewBox(vect.Vector_Zero, 2.7, 0.1)))
	////platform.Physics.Body.SetMoment(Inf)
	//platform.Physics.Shape.SetFriction(0)
	//platform.Physics.Shape.SetElasticity(0)
	//platform.Transform().SetPositionf(1.05, 0.95)
	////platform.Transform().SetScalef(1, 0.1)
	//g.AddGameObject(platform)

	//platform = NewGameObject("platform 2")
	//platform.AddComponent(NewPhysicsShape(true, g.Space, chipmunk.NewBox(vect.Vector_Zero, 2.7, 0.1)))
	////platform.Physics.Body.SetMoment(Inf)
	//platform.Physics.Shape.SetFriction(0)
	//platform.Physics.Shape.SetElasticity(0)
	//platform.Transform().SetPositionf(-1.75, 1.75)
	////platform.Transform().SetScalef(1, 0.1)
	//g.AddGameObject(platform)
	//platform = NewGameObject("platform 3")
	//platform.AddComponent(NewPhysicsShape(true, g.Space, chipmunk.NewBox(vect.Vector_Zero, 2.7, 0.1)))
	////platform.Physics.Body.SetMoment(Inf)
	//platform.Physics.Shape.SetFriction(0)
	//platform.Physics.Shape.SetElasticity(0)
	//platform.Transform().SetPositionf(1.75, 2.15)
	////platform.Transform().SetScalef(1, 0.1)
	//g.AddGameObject(platform)

	answer := PacketResultJoinGame{id, 0}
	answer.Write(c.TCPWriter)
}
