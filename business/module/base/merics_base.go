package base

type MetricsBase struct {
	Name string
}

func (m *MetricsBase) GetName() string {
	return m.Name
}

func (m *MetricsBase) SetName(str string) {
	m.Name = str
}
