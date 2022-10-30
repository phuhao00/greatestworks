package npc

type Abstract interface {
	LoadCfgTableData()
	LoadBehaviorTree(treeFile string)
	SendToRangePlayers(msgId uint64, msg interface{})
	SendUpdateInfoToRangePlayers()
	Update()
	MoveOneStep()
	DoPatrol()
	DoFollow(target interface{})
	PathLength() float32
	IsArrivePoint()
	SetRandomStatus()
	SetRandomAction()
	SwitchingRandomAction()
	DoOnceAction()
	RandContinuousAction()
	SetRandomPatrolPoints()
	SpecialActionCondition()
	AppearCondition()
	DisappearCondition()
	EveryDayConditionCheck()
	EveryWeekConditionCheck()
	AtSpecificTimeCondition()
	WeatherConditionCheck()
	PortalConditionCheck()
	AppearTrigger()
	DisappearTrigger()
	CheckPatrolTime()
	MoveToTarget()
	Appear()
	Disappear()
}

type Aoi interface {
	CLX()
	CLZ()
	CLRange()
	CLEntityID()
	CLEntityType()
	OnEnterRange()
	OnLeaveRange()
}
