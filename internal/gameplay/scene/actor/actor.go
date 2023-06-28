package actor

type Actor interface {
	OnDamage(delta int64)
	Attack()
	OnMove()
}
