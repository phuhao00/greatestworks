package pet

type Abstract interface {
}

type NormalAction interface {
	Eat()
	Move()
}

type FightAction interface {
	Fight()
}
