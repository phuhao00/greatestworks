package internal

// MetricsBase  ...
type MetricsBase struct {
	Name string
}

func (m *MetricsBase) Description() string {
	//TODO implement me
	panic("implement me")
}

func (m *MetricsBase) GetName() string {
	return m.Name
}

func (m *MetricsBase) SetName(str string) {
	m.Name = str
}
