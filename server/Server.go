package server

import (
	"bufio"
	//"bytes"
	"io"
	"log"
	"net"
	"strings"
)

var (
	MainServer *Server
	//EndPacket  []byte = []byte{128, 0, 128, 1}
)

type Job func()

type PacketHandler func(c *Client, rawPacket Packet)
type PacketProcessor func([]byte) Packet

type Server struct {
	Socket        *net.TCPListener
	Clients       map[ID]*Client
	Jobs          chan Job
	IDGen         *IDGenerator
	PacketHandle  []PacketHandler
	PacketProcess []PacketProcessor
}

func (s *Server) Run() {
	for job := range s.Jobs {
		job()
	}
}

type Client struct {
	Socket  *net.TCPConn
	UDPCon  *net.UDPConn
	UDPAddr *net.UDPAddr

	UDPWriter *bufio.Writer
	TCPReader *bufio.Reader
	TCPWriter *bufio.Writer

	ID       ID
	Name     string
	X, Y     float32
	Rotation float32
	Game     *Game_t
	GameID   uint32

	Character    *CharacterController
	Disconnected int32
}

func isTransportOver(data string) (over bool) {
	over = strings.HasSuffix(data, "\r\n\r\n")
	return
}

func (c *Client) Update() {}

func (c *Client) Run() {
	defer c.OnPanic()
	log.Println("Income connection")

	c.TCPReader = bufio.NewReader(c.Socket)
	c.TCPWriter = bufio.NewWriter(c.Socket)

	//c.UDPCon = net.ResolveUDPAddr(c.Socket.RemoteAddr().Network(), c.Socket.RemoteAddr().String())

	var (
		buf = make([]byte, 1024)
	)

	for {
		bufLength, err := c.TCPReader.Read(buf) //[0:])

		//log.Printf("buffer length = %d", bufLength)

		switch err {

		case io.EOF:
			return
		case nil:
			{
				var i uint16 = 0
				for true {
					pLength := uint16(buf[i+1]) | (uint16(buf[i+2]) << 8)
					//log.Printf("packet length = %d", pLength)

					copySlice := buf[i : i+pLength] // copy value hope this is work :)
					MainServer.Jobs <- func() { c.HandlePacket(copySlice, pLength) }

					i += pLength
					if i >= (uint16)(bufLength) {
						break
					}
				}
				//for _, val := range packets {
				//	if val[0] != 0 {
				//		lengh := len(val) // copy value
				//		copySlice := val  // copy value hope this is work :)
				//		MainServer.Jobs <- func() { c.HandlePacket(copySlice, lengh) }
				//	}
				//}
				//buf = make([]byte, 512) // clear
			}
		default: // something wrong, when connection was losted
			{
				// (10054 on windows WSAECONNRESET)
				if c.Game != nil {
					if c.Character != nil {
						var id uint32 = c.Character.GameObject().GOID
						for e := c.Game.Packages.Front(); e != nil; e = e.Next() {
							if id == e.Value.(PacketResultInstanceGO).GOID {
								c.Game.Packages.Remove(e)
								break
							}
						}
						DestroyCharacter := PacketCharacterState{id, CS_Destroy}
						c.Game.Broadcast2OtherTCP(DestroyCharacter, c.ID, false)
						c.Character.GameObject().Destroy()
						c.Game.RemovePlayer(c)
					}
				}
				delete(MainServer.Clients, c.ID)
				log.Printf(err.Error())
				return
			}
		}
	}
}

//func (c *Client) Run() {
//	defer c.OnPanic()
//	log.Println("Income connection")

//	c.TCPReader = bufio.NewReader(c.Socket)
//	c.TCPWriter = bufio.NewWriter(c.Socket)

//	//c.UDPCon = net.ResolveUDPAddr(c.Socket.RemoteAddr().Network(), c.Socket.RemoteAddr().String())

//	var (
//		buf = make([]byte, 512)
//	)

//	for {
//		//n, err := c.TCPReader.Read(buf)
//		//data := string(buf[:n])
//		_, err := c.TCPReader.Read(buf) //[0:])
//		//n++
//		packets := bytes.Split(buf, EndPacket)

//		switch err {

//		case io.EOF:
//			return
//		case nil:
//			{
//				for _, val := range packets {
//					if val[0] != 0 {
//						lengh := len(val) // copy value
//						copySlice := val  // copy value hope this is work :)
//						MainServer.Jobs <- func() { c.HandlePacket(copySlice, lengh) }
//					}
//				}
//				buf = make([]byte, 512) // clear
//			}
//		default: // something wrong, when connection was losted
//			{
//				// (10054 on windows WSAECONNRESET)
//				if c.Game != nil {
//					if c.Character != nil {
//						var id uint32 = c.Character.GameObject().GOID
//						for e := c.Game.Packages.Front(); e != nil; e = e.Next() {
//							if id == e.Value.(PacketResultInstanceGO).GOID {
//								c.Game.Packages.Remove(e)
//								break
//							}
//						}
//						DestroyCharacter := PacketCharacterState{id, CS_Destroy}
//						c.Game.Broadcast2OtherTCP(DestroyCharacter, c.ID, false)
//						c.Character.GameObject().Destroy()
//						c.Game.RemovePlayer(c)
//					}
//				}
//				delete(MainServer.Clients, c.ID)
//				log.Printf(err.Error())
//				return
//			}
//		}
//	}
//}

func (c *Client) Send(p Packet) {}

func (c *Client) OnPanic() {
	if x := recover(); x != nil {
		//if atomic.CompareAndSwapInt32(&c.Disconnected, 0, 1) {
		//	log.Println(c.Name, "Disconnected. Reason:", x)
		//	MainServer.Jobs <- func() {
		//		delete(MainServer.Clients, c.ID)
		//		MainServer.IDGen.PutID(c.ID)
		//	}
		//}
	}
}

func (c *Client) HandlePacket(data []byte, lenght uint16) {
	defer c.OnPanic()

	//if MainServer.PacketProcess[data[0]] != nil {
	packet := MainServer.PacketProcess[data[0]](data[:lenght])

	//	if MainServer.PacketHandle[data[0]] != nil {
	MainServer.PacketHandle[data[0]](c, packet)
	//	}
	//}

	//switch PacketID(data[0]) {
	//case ID_Login:
	//	{
	//		packet = ProcessPacketLogin(data[:lenght])
	//		HandlePacketLogin(c, packet)
	//	}
	//case ID_PlayerInput:
	//	{
	//		packet = HandlePacketPlayerInput(data[:lenght])
	//		OnPacketPlayerInput(c, packet.(PacketPlayerInput))
	//	}
	//case ID_RequestGames:
	//	{
	//		packet = HandlePacketRequestGames(data[:lenght])
	//		OnPacketRequestGames(c, packet.(PacketRequestGames))
	//	}
	//case ID_CreateGame:
	//	{
	//		packet = HandlePacketCreateGame(data[:lenght])
	//		OnPacketGameCreate(c, packet.(PacketGameCreate))
	//	}
	//case ID_JoinGame:
	//	{
	//		packet = HandlePacketJoinGame(data[:lenght])
	//		OnPacketJoinGame(c, packet.(PacketJoinGame))
	//	}
	////case ID_ResolveUDP:
	////	{
	////		packet = HandlePacketResolveUPD(data[:lenght])
	////		OnPacketResolveUPD(c, packet.(PacketResolveUPD))
	////	}
	//case ID_InstanceGO:
	//	{
	//		packet = HandlePacketInstanceGO(data[:lenght])
	//		OnPacketInstanceGO(c, packet.(PacketInstanceGO))
	//	}
	//case 60:
	//	{
	//		log.Printf("packet: id=%d len=%d", data[0], lenght)
	//		var str string = "<cross-domain-policy><allow-access-from domain=\"*\" to-ports=\"*\"/></cross-domain-policy>"
	//		c.TCPWriter.WriteString(str)
	//		c.TCPWriter.Flush()
	//	}
	//default:
	//	{
	//		log.Printf("Unhandled packet: id=%d len=%d", data[0], lenght)
	//	}
	//}
}

func AcceptUDP(UDP_Listner *net.UDPConn) {
	for {
		var (
			buf      = make([]byte, 1024)
			PlayerID ID
		)

		_, addr, err := UDP_Listner.ReadFromUDP(buf[0:])

		if err != nil {
			log.Printf("AcceptUDP error:" + err.Error())
			continue
		}
		if buf[0] == ID_ResolveUDP {
			PlayerID = ID(buf[3]) | (ID(buf[4]) << 8) | (ID(buf[5])<<16 | ID(buf[6])<<24)
			for _, c := range MainServer.Clients {
				if PlayerID == c.ID { // TODO: must be reply TCP message with approve connection
					log.Printf("%s pid=%d", addr.String(), PlayerID)
					c.UDPCon = UDP_Listner
					c.UDPAddr = addr
				}
			}
			buf = make([]byte, 1024)
			continue
		}
	}
}

func StartServer() {
	TCP_addr, TCP_err := net.ResolveTCPAddr("tcp", "0.0.0.0:4354")
	if TCP_err != nil {
		log.Println(TCP_err)
		return
	}
	ln, err := net.ListenTCP("tcp", TCP_addr)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Server started (TCP)! at [%s]", TCP_addr)

	UDP_addr, UDP_err := net.ResolveUDPAddr("udp4", "0.0.0.0:4354")
	if UDP_err != nil {
		log.Println(UDP_err)
		return
	}

	UDP_Listner, err := net.ListenUDP("udp", UDP_addr)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Server started (UDP)! at [%s]", UDP_addr)
	//MainServer.IDGen can be not safe because the only place we use it is when
	//we adding/removing clients from the list and we need to do it safe anyway
	//Socket  *net.TCPListener	// -- ln, err := net.ListenTCP("tcp", TCP_addr)
	//Clients map[ID]*Client		// -- make(map[ID]*Client)
	//Jobs    chan Job			// -- make(chan Job, 1000)
	//IDGen   *IDGenerator		// -- NewIDGenerator(100000, false)
	MainServer = &Server{ln, make(map[ID]*Client), make(chan Job, 1000), NewIDGenerator(100000, false), make([]PacketHandler, ID_Count), make([]PacketProcessor, ID_Count)}

	// [login]
	MainServer.PacketProcess[ID_Login] = ProcessPacketLogin
	MainServer.PacketHandle[ID_Login] = HandlePacketLogin

	// [InstanceGO]
	MainServer.PacketProcess[ID_InstanceGO] = ProcessPacketInstanceGO //
	MainServer.PacketHandle[ID_InstanceGO] = HandlePacketInstanceGO   //

	// [JoinGame]
	MainServer.PacketProcess[ID_JoinGame] = ProcessPacketJoinGame //
	MainServer.PacketHandle[ID_JoinGame] = HandlePacketJoinGame   //

	// [PlayerInput]
	MainServer.PacketProcess[ID_PlayerInput] = ProcessPacketPlayerInput //
	MainServer.PacketHandle[ID_PlayerInput] = HandlePacketPlayerInput   //

	// [RequestGames]
	MainServer.PacketProcess[ID_RequestGames] = ProcessPacketRequestGames //
	MainServer.PacketHandle[ID_RequestGames] = HandlePacketRequestGames   //

	// [CreateGame]
	MainServer.PacketProcess[ID_CreateGame] = ProcessPacketCreateGame
	MainServer.PacketHandle[ID_CreateGame] = HandlePacketCreateGame

	// [ResolveUDP <- login]

	//packet = HandlePacketPlayerInput(data[:lenght])
	//OnPacketPlayerInput(c, packet.(PacketPlayerInput))
	//packet = HandlePacketRequestGames(data[:lenght])
	//OnPacketRequestGames(c, packet.(PacketRequestGames))
	//packet = HandlePacketCreateGame(data[:lenght])
	//HandlePacketGameCreate(c, packet.(PacketGameCreate))
	//packet = HandlePacketJoinGame(data[:lenght])
	//OnPacketJoinGame(c, packet.(PacketJoinGame))
	//packet = HandlePacketInstanceGO(data[:lenght])
	//OnPacketInstanceGO(c, packet.(PacketInstanceGO))

	//log.Printf("packet: id=%d len=%d", data[0], lenght) // [60] packet unity3d web player
	//var str string = "<cross-domain-policy><allow-access-from domain=\"*\" to-ports=\"*\"/></cross-domain-policy>"
	//c.TCPWriter.WriteString(str)
	//c.TCPWriter.Flush()

	go MainServer.Run()
	//go MainServer.RunGameLoops()

	go AcceptUDP(UDP_Listner)

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			log.Println(err)
			break
		}
		MainServer.Jobs <- func() {
			id := MainServer.IDGen.NextID()
			c := &Client{Socket: conn, ID: id}
			MainServer.Clients[c.ID] = c
			go c.Run()
		}
	}
}
