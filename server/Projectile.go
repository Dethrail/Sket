package server

import (
	"log"
)

type Projectile struct {
	BaseComponent
	Source   *CharacterController
	Physics  *Physics
	LifeTime float32
}

func (p *Projectile) OnCollisionEnter(arbiter Arbiter) bool {
	if arbiter.GameObjectB() != nil && arbiter.GameObjectB() != p.Source.GameObject() {
		var cc *CharacterController = nil
		enemy := arbiter.GameObjectB()
		cc, _ = enemy.ComponentTypeOf(cc).(*CharacterController)
		if cc != nil {
			cc.HitPoints = 0
			arbiter.GameObjectA().Destroy()
		}
	}
	return true
}

func (p *Projectile) OnCollisionExit(arbiter Arbiter) {
	if arbiter.GameObjectB() != nil && arbiter.GameObjectB() != p.Source.GameObject() {
		log.Printf("GO exit name = %d, layer=%d", arbiter.GameObjectB().GOID, arbiter.GameObjectB().Physics.Shape.Group)
	}
}

func NewProjectile(c *CharacterController, p *Physics, lifeTime float32) *Projectile {
	return &Projectile{NewComponent(), c, p, lifeTime}
}

func (p *Projectile) Start() {
	//p.Physics.Body.SetMass(1)
	//p.Physics.Shape.Group = 1
}

func (p *Projectile) FixedUpdate() {
	p.LifeTime -= FixedTimeStep
	if p.LifeTime < 0 {
		//log.Printf("destroy projectile")
		p.GameObject().Destroy()
	}
}

func (p *Projectile) Update() {
	//log.Printf("%f, %f", p.Transform().Position().X, p.Transform().Position().Y)
}

func (p *Projectile) LateUpdate() {}

func (p *Projectile) OnDestroy() {
	//log.Printf("OnDestroy")
}
