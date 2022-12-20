package config

type Develop struct {
	Mode      string `yaml:"mode"`
	LogFolder string `yaml:"logFolder"`
}

type Mode string

const (
	DevelopMode Mode = "dev"
	QAMode      Mode = "qa"
	ReleaseMode Mode = "release"
)
