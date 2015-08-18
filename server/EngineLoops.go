package server

import ()

func onCollisionPreSolveGameObject(gameObject *GameObject, arb Arbiter) bool {
	if gameObject == nil || !gameObject.active || arb.GameObjectB() == nil {
		return true
	}
	l := len(gameObject.components)
	comps := gameObject.components

	b := true
	for i := l - 1; i >= 0; i-- {
		b = b && comps[i].OnCollisionPreSolve(arb)
	}
	return b
}

func onCollisionPostSolveGameObject(gameObject *GameObject, arb Arbiter) {
	if gameObject == nil || !gameObject.active || arb.GameObjectB() == nil {
		return
	}
	l := len(gameObject.components)
	comps := gameObject.components

	for i := l - 1; i >= 0; i-- {
		if comps[i].Started() {
			comps[i].OnCollisionPostSolve(arb)
		}
	}
}

func onCollisionEnterGameObject(gameObject *GameObject, arb Arbiter) bool {
	if gameObject == nil || !gameObject.active || arb.GameObjectB() == nil {
		return true
	}
	l := len(gameObject.components)
	comps := gameObject.components

	b := true
	for i := l - 1; i >= 0; i-- {
		if comps[i].Started() {
			b = b && comps[i].OnCollisionEnter(arb)
		}
	}
	return b
}

func onCollisionExitGameObject(gameObject *GameObject, arb Arbiter) {
	if gameObject == nil || !gameObject.active || arb.GameObjectB() == nil {
		return
	}
	l := len(gameObject.components)
	comps := gameObject.components

	for i := l - 1; i >= 0; i-- {
		if comps[i].Started() {
			comps[i].OnCollisionExit(arb)
		}
	}
}
