package server

import (
	//"bytes"
	"fmt"

	"github.com/Dethrail/chipmunk"
	"github.com/Dethrail/chipmunk/vect"
	"log"
)

func (p PacketInstanceGO) ID() PacketID {
	return ID_InstanceGO
}

type PacketInstanceGO struct {
	GameObjectType uint32
}

func ProcessPacketInstanceGO(data []byte) Packet {
	var (
		goType uint32 = 0
	)

	//rawData := bytes.Split(data, []byte{129})
	//for i, val := range rawData {
	//	switch i {
	//	case 0:
	//		{
	id := data[0]
	if id != ID_InstanceGO {
		log.Printf("id of packet (%d) != ID_InstanceGO(%d)", id, ID_InstanceGO)
	}

	//packetLength := uint16(data[1]) | (uint16(data[2]) << 8)
	goType = uint32(data[3]) | (uint32(data[4]) << 8) | (uint32(data[5])<<16 | uint32(data[6])<<24)
	//log.Printf("gameObject type(%d)", goType)
	//}

	//default:
	//{
	//	log.Printf("intval %d = %d", i, val[0])
	//}
	//}
	//}

	//for i, val := range arr {

	return PacketInstanceGO{goType}
}

func HandlePacketInstanceGO(c *Client, rawPacket Packet) {
	var (
		p      PacketInstanceGO = rawPacket.(PacketInstanceGO)
		id     uint32           = 0
		result byte             = 1
		x      float32          = 0
		y      float32          = 0
	)
	if p.GameObjectType == 0 && Games[c.GameID] != nil {
		pName := fmt.Sprintf("Player(%i)", c.ID)
		character := NewGameObject(pName)
		//physics := character.AddComponent( NewPhysicsShape(false, Games[c.GameID].Space, chipmunk.NewCircle(vect.Vector_Zero, float32(0.06)))).(* Physics) //chipmunk.NewBox(vect.Vector_Zero, 0.12, 0.12)))
		physics := character.AddComponent(NewPhysicsShape(false, Games[c.GameID].Space, chipmunk.NewBox(vect.Vector_Zero, 0.04, 0.12))).(*Physics) //
		cc := character.AddComponent(NewCharacterController(Games[c.GameID], c.ID)).(*CharacterController)
		c.Character = cc
		physics.Body.SetMoment(Inf)
		physics.Shape.SetFriction(1)
		physics.Shape.SetElasticity(0)
		character.Transform().SetPositionf(0, 1)
		x = character.Transform().Position().X
		y = character.Transform().Position().Y
		Games[c.GameID].AddGameObject(character)
		id = character.GOID
		result = 0
	}
	answer := PacketResultInstanceGO{id, uint32(c.ID), x, y, result}
	Games[c.GameID].Broadcast2AllPlayersTCP(answer)

	if result == 0 {
		Games[c.GameID].Packages.PushBack(answer)
		//var i int = 0
		//for e := Games[c.GameID].Packages.Front(); e != nil; e = e.Next() {
		//i++
		//}
		//log.Panicf("Count %d", i)
	}
}

//func (p PacketResolveUPD) Write(w *bufio.Writer) {
//	w.WriteByte(byte(p.ID()))

//	err := w.Flush()

//	if err != nil {
//		log.Print("Error on game obj position" + err.Error() + "\n")
//	}
//}
