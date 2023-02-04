package excel_object

type Object interface {
	Check()
	Patch()
}

type ObjectManager interface {
	Get(key any) Object
	Load(path string) error //加载配置
	LoadAfter() error
}
