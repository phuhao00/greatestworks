package pet

import "greatestworks/business/module/skill"

type Pet struct {
	Category   int //类型
	ConfigId   uint32
	Id         uint64
	Name       string
	Star       uint32
	Level      uint32
	State      State //状态
	Skills     []skill.Skill
	Property   map[uint32]int64 //属性
	Trammels   *Trammels        //羁绊
	ReliveTime int64            //复活时间
}

func (p *Pet) ToModel() *Model {
	return &Model{
		Category: p.Category,
		Id:       p.Id,
		ConfigId: p.ConfigId,
		Name:     p.Name,
	}
}

func (p *Pet) AddLevel(delta uint32) {

}
