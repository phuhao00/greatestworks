package template

type Equip struct {
	*ItemBase
}

func (e *Equip) DayGetCheck() bool {
	return true
}

func (e *Equip) WeekGetCheck() bool {
	return true
}

func (e *Equip) UseCdCheck() bool {
	return true
}
