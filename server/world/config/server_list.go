package config

type ServerList struct {
	ProcIndex int
	MaxPlayer int32
	Name      string
	Name1     string
}

func GetServerList() *ServerList {
	return nil
}
