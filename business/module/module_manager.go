package module

var (
	MManager Manager
)

type Manager struct {
	moduleName2Module map[string]IModule
}

func (m *Manager) GetModule(name string) IModule {
	return m.moduleName2Module[name]
}
