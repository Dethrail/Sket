package server

import (
//"bufio"
//"bytes"
//"fmt"
//"log"
//"net"
)

func (p PacketResolveUPD) ID() PacketID {
	return ID_ResolveUDP
}

type PacketResolveUPD struct {
	Port  uint32
	Error byte
}

//func HandlePacketResolveUPD(data []byte) PacketResolveUPD {
//	//var (
//	//	Port  uint32 = 0
//	//	Error byte   = 0
//	//)

//	////rawData := bytes.Split(data, []byte{129})
//	//for i, val := range rawData {
//	//	switch i {
//	//	case 0:
//	//		{
//	//			id := val[0]
//	//			if id != ID_ResolveUDP {
//	//				log.Printf("id of  packet (%d) != ID_ResolveUDP(%d)", id, ID_ResolveUDP)
//	//			}
//	//		}
//	//	case 1:
//	//		{
//	//			Port = uint32(val[0]) | (uint32(val[1]) << 8) | (uint32(val[2])<<16 | uint32(val[3])<<24)
//	//			//log.Printf("Receive port(%d)", Port)
//	//		}

//	//	default:
//	//		{
//	//			log.Printf("resolve UDP unhandled data intval %d = %d", i, val[0])
//	//		}
//	//	}
//	//}
//	//return PacketResolveUPD{Port, Error}
//}

//func OnPacketResolveUPD(c *Client, p PacketResolveUPD) {
//	//ip, _, err := net.SplitHostPort(c.Socket.RemoteAddr().String())
//	////if err == nil {
//	////	log.Println("ip=%d, port=%d", ip, port) // [remote=192.168.1.1:10726] [local=192.168.1.55 4354]
//	////}
//	//if err == nil {
//	//	str := fmt.Sprintf("%s:%d", ip, p.Port)
//	//	serverAddr, _ := net.ResolveUDPAddr("udp4", str) // _ = error
//	//	log.Printf(str + "test")
//	//	c.UDPCon, err = net.DialUDP("udp", nil, serverAddr)
//	//	c.UDPWriter = bufio.NewWriter(c.UDPCon)
//	//	if err != nil {
//	//		log.Print("error " + err.Error())
//	//	}
//	//}
//}

//func (p PacketResolveUPD) Write(w *bufio.Writer) {
//	//w.WriteByte(byte(p.ID()))

//	//err := w.Flush()

//	//if err != nil {
//	//	log.Print("Error on game obj position" + err.Error() + "\n")
//	//}
//}
