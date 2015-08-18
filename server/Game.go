package server

import (
	"container/list"
	"github.com/Dethrail/chipmunk"
	"github.com/Dethrail/chipmunk/vect"
	"log"
	"math"
	"runtime"
	"time"
)

const (
	RadianConst       = math.Pi / 180
	DegreeConst       = 180 / math.Pi
	MaxPhysicsTime    = float64(1) / float64(60)
	stepTime          = float64(1) / float64(60)
	BroadcastStepRate = 4 // one of four packet will be sanded
)

var (
	Inf   = vect.Float(math.Inf(1))
	Games map[uint32]*Game_t

	GameCounter uint32 = 0

	CorrectWrongPhysics = true
	CountObjects        = 0

	ServerTime float32 = 0
)

type Game_t struct {
	ID          uint32
	Players     []*Client
	GameObjects []*GameObject

	Space     *chipmunk.Space
	deltaTime float64
	fixedTime float64
	gameTime  time.Time
	lastTime  time.Time

	running        bool
	insideGameloop bool
	Count          int
	BroadcastStep  byte
	Packages       *list.List
}

func (game *Game_t) GetGame() (g *Game_t) {
	g = game
	return g
}

func InitGames() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.LockOSThread()
	Games = make(map[uint32]*Game_t) //&server{ln, make(map[int32]* Client), make(chan Job, 1000), NewIDGenerator(100000, false)}

	go GamesLoop()
}

func GamesLoop() {
	//for _ = range time.Tick(time.Second / 30) {
	for _ = range time.Tick(time.Second / 60) {
		for i, _ := range Games {
			Games[i].Loop()
		}
		ServerTime += 1
		if ServerTime > 1000 {
			ServerTime = 0
			runtime.GC()
		}
	}
}

func CreateGame(client *Client) (g *Game_t, id uint32) {
	id = GameCounter
	GameCounter += 1
	g = &Game_t{id, make([]*Client, 0), make([]*GameObject, 0), nil, 0, 0, time.Now(), time.Now(), true, false, -1, 0, list.New()}
	g.RegisterPlayer(client)

	g.Space = chipmunk.NewSpace()
	g.Space.Gravity = vect.Vect{0, -9.8}

	RegisterGame(g)
	return g, id
}

func (g *Game_t) RegisterPlayer(c *Client) {
	g.Players = append(g.Players, c)
	c.Game = g
	c.GameID = g.ID
}
func (g *Game_t) RemovePlayer(client *Client) {
	for i, player := range g.Players {
		if player == client {
			g.Players[i] = nil
			break
		}
	}
	for i := 0; i < len(g.Players); i++ {
		if g.Players[i] == nil {
			g.Players[i], g.Players = g.Players[len(g.Players)-1], g.Players[:len(g.Players)-1]
			i--
		}
	}
	client.Game = nil
	client.GameID = 0
}

func RegisterGame(g *Game_t) {
	Games[g.ID] = g
}

func (g *Game_t) AddGameObject(gameObject ...*GameObject) {
	for _, obj := range gameObject {
		g.GameObjects = append(g.GameObjects, obj)
	}
}

func (g *Game_t) removeGameObject(gameObject *GameObject) {
	if gameObject == nil {
		return
	}
	for i, c := range g.GameObjects {
		if gameObject == c {
			g.GameObjects[i].transform.childOfScene = false
			for _, t := range gameObject.transform.children {
				g.removeGameObject(t.gameObject)
			}

			g.GameObjects[i] = nil // remove from go's
			gameObject.destroy()   // destroy
			break
		}
	}
}

func (g *Game_t) cleanNil() {
	var flag bool = false
	for i := 0; i < len(g.GameObjects); i++ {
		if g.GameObjects[i] == nil {
			g.GameObjects[i], g.GameObjects = g.GameObjects[len(g.GameObjects)-1], g.GameObjects[:len(g.GameObjects)-1]
			i--
			flag = true
		}
	}
	if flag {
		runtime.GC()
	}
}

func (g *Game_t) RemoveGameObject(gameObject *GameObject) {
	if gameObject == nil {
		return
	}
	gameObject.transform.removeFromParent()
	g.removeGameObject(gameObject)
}

func (g *Game_t) Broadcast2AllPlayersTCP(packet Packet) {
	for _, p := range g.Players {
		w := p.TCPWriter
		//log.Printf("Send p=%d and target=%d", packet.ID(), p.ID)
		packet.Write(w)
	}
}

func (g *Game_t) Broadcast2OtherTCP(packet Packet, packetOwner ID, write2Owner bool) {
	for _, p := range g.Players {
		if p.ID == packetOwner && !write2Owner {
			continue
		}
		w := p.TCPWriter
		packet.Write(w)
	}
}

func (g *Game_t) Loop() bool {
	g.insideGameloop = true
	if g.running {
		g.Run()
	} else {
		return false
	}
	g.insideGameloop = false

	return true
}

func (g *Game_t) Run() {
	//if CountObjects != len(g.GameObjects) {
	//	CountObjects = len(g.GameObjects)
	//	log.Printf("Count = %d", CountObjects)
	//}
	//timeNow := time.Now()
	//gameTime = gameTime.Add(timeNow.Sub(lastTime))
	//g.deltaTime = float64(timeNow.Sub(g.lastTime).Nanoseconds()) / float64(time.Second)
	//g.lastTime = timeNow

	//timer := NewTimer()
	//timer.Start()

	//var destroyDelta,
	//	startDelta time.Duration
	// 	fixedUpdateDelta,
	// 	physicsDelta,
	// 	updateDelta,
	// 	lateUpdateDelta,
	// 	drawDelta,
	// 	coroutinesDelta,
	//var stepDelta time.Duration
	// 	behaviorDelta,
	// 	startPhysicsDelta,
	// 	endPhysicsDelta time.Duration

	//d := g.deltaTime

	//if d > MaxPhysicsTime {
	//	d = MaxPhysicsTime
	//}
	//g.fixedTime += d

	arr := &g.GameObjects

	//timer.StartCustom("Destory routines")
	iterWithGame(arr, g, destoyGameObject)
	//destroyDelta = timer.StopCustom("Destory routines")

	g.cleanNil()

	//timer.StartCustom("Start routines")
	iter(arr, startGameObject)
	//startDelta = timer.StopCustom("Start routines")

	//timer.StartCustom("Physics time")

	//timer.StartCustom("Physics step time")
	//for g.fixedTime >= stepTime {
	//timer.StartCustom("FixedUpdate routines")
	iterWithGame(arr, g, fixedUdpateGameObject)
	g.BroadcastStep++
	if g.BroadcastStep >= BroadcastStepRate {
		g.BroadcastStep = 0
		iterWithGame(arr, g, broadcastObjectsState)
	}
	//fixedUpdateDelta = timer.StopCustom("FixedUpdate routines")

	//timer.StartCustom("PreStep Physics Delta")
	iter(arr, preStepGameObject)
	//startPhysicsDelta
	//_ = timer.StopCustom("PreStep Physics Delta")
	g.Space.Step(vect.Float(stepTime))
	//g.fixedTime -= stepTime

	//timer.StopCustom("Physics step time")

	//timer.StartCustom("PostStep Physics Delta")
	iter(arr, postStepGameObject)
	//endPhysicsDelta = timer.StopCustom("PostStep Physics Delta")
	//}
	//if g.fixedTime > 0 && g.fixedTime < stepTime {
	//	//iterWithGame(arr, g, interpolateGameObject)
	//}
	//log.Print("in loop\n")
	// TODO: broadcast

	//physicsDelta = timer.StopCustom("Physics time")

	//timer.StartCustom("Update routines")
	iter(arr, udpateGameObject)
	//updateDelta = timer.StopCustom("Update routines")

	//timer.StartCustom("LateUpdate routines")
	iter(arr, lateudpateGameObject)
	//lateUpdateDelta = timer.StopCustom("LateUpdate routines")

	//stepDelta = timer.Stop()

	//if stepDelta > 17*time.Millisecond {
	//	log.Println("StepDelta time is lower than normal")
	//}
	if len(g.Players) == 0 {
		log.Printf("Close game")
		delete(Games, g.ID) // close game
	}
}

//func iter(objsp []*GameObject, f func(*GameObject)) {
//	for i := len(objsp) - 1; i >= 0; i-- {
//		if objsp[i] != nil {
//			f(objsp[i])
//		}
//	}
//}

//func iterWithGame(objsp []*GameObject, g *Game_t, f func(*Game_t, *GameObject)) {
//	for i := len(objsp) - 1; i >= 0; i-- {
//		if objsp[i] != nil {
//			f(g, objsp[i])
//		}
//	}
//}

func iter(objsp *[]*GameObject, f func(*GameObject)) {
	objs := *objsp
	for i := len(objs) - 1; i >= 0; i-- {
		obj := objs[i]
		if obj != nil {
			f(obj)
			//Checks if the objs array has been changed
			if i >= len(*objsp) {
				break
			}
			if obj != objs[i] {
				i++
			}
		}
	}
}

func iterWithGame(objsp *[]*GameObject, g *Game_t, f func(*Game_t, *GameObject)) {
	objs := *objsp
	for i := len(objs) - 1; i >= 0; i-- {
		obj := objs[i]
		if obj != nil {
			f(g, obj)
			//Checks if the objs array has been changed
			if i >= len(*objsp) {
				break
			}
			if obj != objs[i] {
				i++
			}
		}
	}
}
