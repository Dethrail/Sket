package server

import (
	"github.com/Dethrail/chipmunk"
	"github.com/Dethrail/chipmunk/vect"
	//"log"
)

type Physics struct {
	BaseComponent
	Space *chipmunk.Space
	Body  *chipmunk.Body
	Box   *chipmunk.BoxShape
	Shape *chipmunk.Shape

	LastPosition vect.Vect
	LastAngle    vect.Float

	InterpolatedPosition vect.Vect
	InterpolatedAngle    vect.Float

	Interpolate bool
}

func NewPhysics(static bool, space *chipmunk.Space) *Physics {
	var body *chipmunk.Body

	box := chipmunk.NewBox(vect.Vect{0, 0}, vect.Float(1), vect.Float(1))

	if static {
		body = chipmunk.NewBodyStatic()
	} else {
		body = chipmunk.NewBody(1, box.Moment(1))
	}

	p := &Physics{BaseComponent: NewComponent(), Space: space, Body: body, Box: box.GetAsBox(), Shape: box}
	body.AddShape(box)
	return p
}

func NewPhysicsCircle(static bool) *Physics {
	var body *chipmunk.Body
	shape := chipmunk.NewCircle(vect.Vector_Zero, 1)
	if static {
		body = chipmunk.NewBodyStatic()
	} else {
		body = chipmunk.NewBody(1, shape.ShapeClass.Moment(1))
	}

	p := &Physics{BaseComponent: NewComponent(), Body: body, Box: shape.GetAsBox(), Shape: shape}
	body.AddShape(shape)
	return p
}

func NewPhysicsShape(static bool, space *chipmunk.Space, shape *chipmunk.Shape) *Physics {

	var body *chipmunk.Body
	if static {
		body = chipmunk.NewBodyStatic()
		//s := chipmunk.NewBox(vect.Vect{-100, -10}, 200, 5)
		//body.AddShape(s)
		//p := &Physics{BaseComponent: NewComponent(), Space: space, Body: body, Box: s.GetAsBox(), Shape: s}
		//return p
	} else {
		body = chipmunk.NewBody(1, shape.ShapeClass.Moment(1))
	}

	body.AddShape(shape)
	p := &Physics{BaseComponent: NewComponent(), Space: space, Body: body, Box: shape.GetAsBox(), Shape: shape}
	return p
}

func NewPhysicsShapes(static bool, shapes []*chipmunk.Shape) *Physics {
	var body *chipmunk.Body
	if static {
		body = chipmunk.NewBodyStatic()
	} else {
		moment := vect.Float(0)
		for _, shape := range shapes {
			moment += shape.Moment(1)
		}
		body = chipmunk.NewBody(1, moment)
	}

	p := &Physics{BaseComponent: NewComponent(), Body: body, Box: nil, Shape: nil}
	for _, shape := range shapes {
		body.AddShape(shape)
	}
	return p
}

func (p *Physics) Start() {
	//p.Interpolate = true
	pos := p.GameObject().Transform().WorldPosition()

	p.Body.SetAngle(vect.Float(p.GameObject().Transform().WorldRotation().Z) * RadianConst)
	p.Body.SetPosition(vect.Vect{vect.Float(pos.X), vect.Float(pos.Y)})
	p.LastPosition = p.Body.Position()
	p.LastAngle = p.Body.Angle()

	//if p.GameObject().Sprite != nil {
	//	p.GameObject().Sprite.UpdateShape()
	//	p.Body.UpdateShapes()
	//}

	p.Body.UpdateShapes()
	//log.Printf("3x=%f, y=%f", p.Body.Shapes[0].BB.Extents().X, p.Body.Shapes[0].BB.Extents().Y)
	p.Space.AddBody(p.Body)
}

func (p *Physics) OnEnable() {
	p.Body.Enabled = true
}

func (p *Physics) OnDisable() {
	p.Body.Enabled = false
}

func (p *Physics) OnComponentAdd() {
	//log.Print(p.gameObject == nil)
	p.gameObject.Physics = p
	p.Body.CallbackHandler = p
}

func (p *Physics) CollisionPreSolve(arbiter *chipmunk.Arbiter) bool {
	if p.gameObject == nil {
		return true
	}
	return onCollisionPreSolveGameObject(p.GameObject(), newArbiter(arbiter, p.gameObject))
}

func (p *Physics) CollisionEnter(arbiter *chipmunk.Arbiter) bool {
	if p.gameObject == nil || arbiter == nil {
		return true
	}
	return onCollisionEnterGameObject(p.GameObject(), newArbiter(arbiter, p.gameObject))
}

func (p *Physics) CollisionExit(arbiter *chipmunk.Arbiter) {
	if p.gameObject == nil {
		return
	}
	onCollisionExitGameObject(p.GameObject(), newArbiter(arbiter, p.gameObject))
}

func (p *Physics) CollisionPostSolve(arbiter *chipmunk.Arbiter) {
	if p.gameObject == nil {
		return
	}
	onCollisionPostSolveGameObject(p.GameObject(), newArbiter(arbiter, p.gameObject))
}

func (p *Physics) OnDestroy() {
	//log.Printf("OnDestroy")
	p.gameObject = nil
	p.Space.RemoveBody(p.Body)
}

func (p *Physics) Clone() {
	p.Body = p.Body.Clone()
	p.Box = p.Body.Shapes[0].GetAsBox()
	p.Shape = p.Body.Shapes[0]
}
