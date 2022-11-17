package family

type Config struct {
}

type MemberPosition int

const (
	Leader MemberPosition = iota + 1
	Normal
)

type MemberAction int

const (
	RemoveMember MemberAction = iota + 1
	RewardGIft
)

type RemoveMemberActionParam struct {
	OpMemberId uint64
	RemoveIds  []uint64
	Reason     string
}

type RewardGiftActionParam struct {
	OpMemberId uint64
	RemoveIds  []uint64
	ItemIds    map[uint32]uint64
}

type MemberActionParam struct {
	Action MemberAction
	Info   interface{}
}

type ManagerHandlerParam struct {
	FamilyId          uint64
	MemberActionParam *MemberActionParam
}
