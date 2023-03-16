package family

func HandleMemberAction(f *Family, param *MemberActionParam) {
	if param.Action == RemoveMember {
		info := param.Info.(*RemoveMemberActionParam)
		member := f.GetMember(info.OpMemberId)
		if member != nil && member.Position == Leader {
			f.DelMember(info.RemoveIds)
			f.ChOut <- struct {
			}{}
		}
	}
}
