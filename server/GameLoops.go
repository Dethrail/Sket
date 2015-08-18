package server

import (
	"github.com/Dethrail/chipmunk/vect"
	//"math"
	//"log"
	"time"
)

func destoyGameObject(game *Game_t, gameObject *GameObject) {
	if gameObject.DestoryMark {
		game.RemoveGameObject(gameObject) // will destroy this object
	}
}

func startGameObject(gameObject *GameObject) {
	if !gameObject.IsActive() {
		return
	}

	comps := gameObject.Components()
	l := len(comps)

	for i := l - 1; i >= 0; i-- {
		if !comps[i].Started() {
			comps[i].SetStarted(true)
			comps[i].Start()
		}
	}
	comps = nil
}

func fixedUdpateGameObject(game *Game_t, gameObject *GameObject) {
	if !gameObject.IsActive() || gameObject.Physics == nil {
		return
	}

	comps := gameObject.Components()
	l := len(comps)

	for i := l - 1; i >= 0; i-- {
		if comps[i].Started() {
			comps[i].FixedUpdate()
		}
	}
	comps = nil
}

func preStepGameObject(g *GameObject) {
	if g.Physics != nil && g.IsActive() && !g.Physics.Body.IsStatic() && g.Physics.Started() {
		pos := g.Transform().WorldPosition()
		angle := g.Transform().Angle() * RadianConst

		if g.Physics.Interpolate {
			//Interpolation check: if position/angle has been changed directly and not by the physics engine, change g.Physics.LastPosition/LastAngle
			if vect.Float(pos.X) != g.Physics.InterpolatedPosition.X || vect.Float(pos.Y) != g.Physics.InterpolatedPosition.Y {
				g.Physics.InterpolatedPosition = vect.Vect{vect.Float(pos.X), vect.Float(pos.Y)}
				g.Physics.Body.SetPosition(g.Physics.InterpolatedPosition)
			}
			if vect.Float(angle) != g.Physics.InterpolatedAngle {
				g.Physics.InterpolatedAngle = vect.Float(angle)
				g.Physics.Body.SetAngle(g.Physics.InterpolatedAngle)
			}
		} else {
			var pPos vect.Vect
			pPos.X, pPos.Y = vect.Float(pos.X), vect.Float(pos.Y)

			g.Physics.Body.SetAngle(vect.Float(angle))
			g.Physics.Body.SetPosition(pPos)
		}
		g.Physics.LastPosition = g.Physics.Body.Position()
		g.Physics.LastAngle = g.Physics.Body.Angle()
	}
}

func postStepGameObject(g *GameObject) {
	if g.Physics != nil && g.IsActive() && !g.Physics.Body.IsStatic() && g.Physics.Started() {
		/*
			When parent changes his position/rotation it changes his children position/rotation too but the physics engine thinks its in different position
			so we need to check how much it changed and apply to the new position/rotation so we wont fuck up things too much.

			Note:If position/angle is changed in between preStep and postStep it will be overrided.
		*/
		if CorrectWrongPhysics {
			b := g.Physics.Body
			angle := float32(b.Angle())

			lAngle := float32(g.Physics.LastAngle)
			lAngle += angle - lAngle

			pos := b.Position()
			lPos := g.Physics.LastPosition
			lPos.X += (pos.X - lPos.X)
			lPos.Y += (pos.Y - lPos.Y)

			if g.Physics.Interpolate {
				g.Physics.InterpolatedAngle = vect.Float(lAngle)
				g.Physics.InterpolatedPosition = lPos
			}

			b.SetPosition(lPos)
			//b.SetAngle(g.Physics.InterpolatedAngle)

			//g.Transform().SetWorldRotationf(lAngle * DegreeConst)
			g.Transform().SetWorldPositionf(float32(pos.X), float32(pos.Y))
		} else {
			//b := g.Physics.Body
			//angle := b.Angle()
			//pos := b.Position()

			//if g.Physics.Interpolate {
			//	g.Physics.InterpolatedAngle = angle
			//	g.Physics.InterpolatedPosition = pos
			//}

			////g.Transform().SetWorldRotationf(float32(angle) * DegreeConst)
			//g.Transform().SetWorldPositionf(float32(pos.X), float32(pos.Y))
		}
	}
}

func interpolateGameObject(game *Game_t, g *GameObject) {
	if g.Physics != nil && g.Physics.Interpolate && g.IsActive() && !g.Physics.Body.IsStatic() && g.Physics.Started() {
		nextPos := g.Physics.Body.Position()
		currPos := g.Physics.LastPosition

		//nextAngle := g.Physics.Body.Angle()
		//currAngle := g.Physics.LastAngle

		alpha := vect.Float(game.fixedTime / stepTime)
		x := currPos.X + ((nextPos.X - currPos.X) * alpha)
		y := currPos.Y + ((nextPos.Y - currPos.Y) * alpha)
		//a := currAngle + ((nextAngle - currAngle) * alpha)
		g.Transform().SetWorldPositionf(float32(x), float32(y))
		//g.Transform().SetWorldRotationf(float32(a) * DegreeConst)

		//g.Physics.InterpolatedAngle = a
		g.Physics.InterpolatedPosition.X, g.Physics.InterpolatedPosition.Y = x, y
	}
}

func udpateGameObject(gameObject *GameObject) {
	if !gameObject.IsActive() {
		return
	}

	comps := gameObject.Components()
	l := len(comps)

	for i := l - 1; i >= 0; i-- {
		if comps[i].Started() {
			comps[i].Update()
		}
	}
	comps = nil
}

func lateudpateGameObject(gameObject *GameObject) {
	if !gameObject.IsActive() {
		return
	}

	comps := gameObject.Components()
	l := len(comps)

	for i := l - 1; i >= 0; i-- {
		if comps[i].Started() {
			comps[i].LateUpdate()
		}
	}
	comps = nil
}

func broadcastObjectsState(game *Game_t, gameObject *GameObject) {
	if !gameObject.IsActive() || gameObject.Physics == nil {
		return
	}

	var x float32 = gameObject.Transform().Position().X
	var y float32 = gameObject.Transform().Position().Y

	var vx float32 = float32(gameObject.Physics.Body.Velocity().X)
	var vy float32 = float32(gameObject.Physics.Body.Velocity().Y)

	var rY uint16 = uint16(gameObject.Transform().Rotation().Y)
	var t int32 = (int32)(time.Now().Nanosecond())
	packet := PacketGOState{gameObject.GOID, x, y, vx, vy, rY, t}

	for _, p := range game.Players {
		//if p == nil { // TODO: Remove this is memory leak
		packet.Write(p.UDPCon, p.UDPAddr)
		//}
	}
}
