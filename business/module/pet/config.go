package pet

type Config struct {
	Id              uint32
	Name            string
	Category        int      //分类，类型
	ComposeFragment []uint32 //合成需要的碎片
}

type State uint32

const (
	Fight State = iota + 1
)
