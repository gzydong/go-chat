package config

type App struct {
	Env     string `json:"env"`
	Debug   bool   `json:"debug"`
	JuheKey string `json:"juhe_key" yaml:"juhe_key"`
}
