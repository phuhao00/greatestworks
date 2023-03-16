package scene

type DamageInfo struct {
	ObjectId uint64
	Damage   int64
}

type Fight struct {
	Base
	DamageArr []DamageInfo
}
