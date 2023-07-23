package scene

import b3 "github.com/magicsea/behavior3go"

type IActor interface {
	OnDamage(delta int64)
	Attack()
	OnMove()
	Patrol() b3.Status
	FollowTarget() IActor
	Follow(IActor) b3.Status
	RandomStatus() b3.Status
	AppearTrigger() b3.Status
	DisappearTrigger() b3.Status
	MoveToTarget() b3.Status
}
