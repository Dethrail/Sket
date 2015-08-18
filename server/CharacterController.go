package server

import (
	"github.com/Dethrail/chipmunk"
	"github.com/Dethrail/chipmunk/vect"
	"math/rand"
)

const (
	MaxInertion              = 5
	MaxShiftInertion         = 8
	MoveSpeed        float32 = 5
	ReviveTimer      int     = 5
	FixedTimeStep    float32 = 0.033
)

type FloatVelocity struct {
	X float32
	Y float32
}

type MoveDirection struct {
	X int8
	Y int8
}

type KeyInput struct {
	Space          byte
	Shift          byte
	lastSpaceState byte // check state for press and pressed

	A byte
	S byte
	D byte
	W byte
}

const (
	WarriorAttackHitTime  = 0.3
	WarriorAttackCooldown = 0.75
)

type CharacterController struct {
	BaseComponent

	PlayerID  ID
	HitPoints uint8
	Speed     float32
	JumpSpeed float32
	Physics   *Physics
	state     int

	deathState  bool
	reviveTimer float32
	inertion    float32

	Direction MoveDirection
	Input     KeyInput
	Attack    bool
	Velocity  FloatVelocity

	Ground *GameObject
	Game   *Game_t
}

func (cc *CharacterController) OnCollisionEnter(arbiter Arbiter) bool {

	if arbiter.GameObjectB().Name() == "ground" || arbiter.GameObjectB().Name()[0:3] == "pla" {
		cc.Ground = arbiter.GameObjectB()
		for i := 0; i < arbiter.Arbiter.NumContacts; i++ {
			if arbiter.Normal(arbiter.Arbiter.Contacts[i]).Y > -0.7 {
				arbiter.Arbiter.Ignore()
			}
		}

	}

	return true
}

func (cc *CharacterController) OnCollisionExit(arbiter Arbiter) {
	if arbiter.GameObjectB() == cc.Ground {
		cc.Ground = nil
	}
}

func NewCharacterController(g *Game_t, pid ID) *CharacterController {
	return &CharacterController{NewComponent(), pid, 100, 0.1, 20, nil, -1, false, 0, 0, MoveDirection{0, 0}, KeyInput{}, false, FloatVelocity{0, 0}, nil, g}
}

func (cc *CharacterController) Start() {
	if cc.GameObject().Physics == nil {
		return
	}
	cc.Physics = cc.GameObject().Physics
	cc.Physics.Body.SetMass(1)
	cc.Physics.Shape.Group = 1
}

func (cc *CharacterController) Revive() {
	cc.deathState = false
	PlayRevive := PacketCharacterState{cc.GameObject().GOID, CS_Revive}
	cc.Game.Broadcast2AllPlayersTCP(PlayRevive)
	cc.HitPoints = 100

	cc.Velocity.X = 0
	cc.Velocity.Y = 0
	cc.Physics.Body.SetVelocityX(0)
	cc.Physics.Body.SetVelocityY(0)
	posX := 4 - rand.Intn(8)
	cc.Transform().SetPositionf(float32(posX), 1.4)

	cc.Transform().SetRotationf(0)
}

func (cc *CharacterController) Update() {
	if cc.HitPoints <= 0 {
		//log.Printf("Character Dead, pid=%d", cc.PlayerID)
		if !cc.deathState {
			cc.deathState = true
			PlayDeath := PacketCharacterState{cc.GameObject().GOID, CS_Death}
			cc.Game.Broadcast2AllPlayersTCP(PlayDeath)
		}

		if cc.reviveTimer < (float32)(ReviveTimer) {
			cc.reviveTimer += FixedTimeStep
			return
		} else {
			cc.reviveTimer = 0
			cc.Revive()
		}
	}

	if cc.Input.Space == 1 && cc.Input.Space != cc.Input.lastSpaceState {
		//broadcast attack state
		PlayAttack := PacketCharacterState{cc.GameObject().GOID, CS_Attack}
		cc.Game.Broadcast2OtherTCP(PlayAttack, cc.PlayerID, false)

		attack := NewGameObject("attack")
		physics := attack.AddComponent(NewPhysicsShape(false, cc.Physics.Space, chipmunk.NewCircle(vect.Vector_Zero, float32(0.15)))).(*Physics)
		physics.Shape.IsSensor = true
		physics.Body.IgnoreGravity = true

		var v Vector = cc.Transform().Position()
		if cc.Transform().Rotation().Y == 0 { // l = 180 r =0
			v.X += 0.06
		} else {
			v.X -= 0.06
		}
		attack.Transform().SetPosition(v)
		attack.AddComponent(NewProjectile(cc, physics, 0.1))
		cc.Game.AddGameObject(attack)
	}

	cc.Input.lastSpaceState = cc.Input.Space

	switch cc.Direction.Y {
	case 1:
		{
			//if cc.Ground != nil && cc.Physics.Body.Velocity().Y <= 0 {
			//	cc.Velocity.Y = 5
			//}
			if cc.Ground != nil && cc.Physics.Body.Velocity().Y == 0 {
				cc.Physics.Body.AddForce(0, 325)
			}
		}
	case -1:
		{
			if cc.Ground != nil && cc.Ground.GameObject().Name()[0:3] == "pla" {

				//log.Printf("s")
			}
		}
	}
	var maxInertion float32 = MaxInertion // MaxShiftInertion
	if cc.Input.Shift == 1 {
		maxInertion = MaxShiftInertion
	}

	switch cc.Direction.X {
	case 1:
		{
			if cc.inertion < maxInertion {
				cc.inertion++
			}

			if cc.Input.Shift == 0 && cc.inertion > maxInertion {
				cc.inertion--
			}
		}
	case -1:
		{
			if cc.inertion > -maxInertion {
				cc.inertion--
			}
			if cc.Input.Shift == 0 && cc.inertion < -maxInertion {
				cc.inertion++
			}
		}

	default:
		{
			if cc.inertion != 0 {
				if cc.inertion > 0 {
					cc.inertion--
				} else {
					cc.inertion++
				}
			}
		}
	}

	cc.Velocity.X = MoveSpeed * cc.inertion * FixedTimeStep

	cc.Physics.Body.SetVelocityX(cc.Velocity.X)
	//cc.Physics.Body.SetVelocity(cc.Velocity.X, cc.Velocity.Y)

	s := cc.Transform().Rotation()
	if cc.Direction.X < 0 && s.Y != 180 {
		s.Y = 180
	} else if cc.Direction.X > 0 && s.Y != 0 {
		s.Y = 0
	}
	cc.Transform().SetRotation(s)
}

func (cc *CharacterController) LateUpdate() {
	//gameSceneGeneral.SceneData.Camera.Transform().SetPosition(NewVector3(300-cc.Transform().Position().X, 0, 0))
}
