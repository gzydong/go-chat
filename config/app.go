package config

type App struct {
	Env        string   `json:"env"`
	Debug      bool     `json:"debug"`
	JuheKey    string   `json:"juhe_key" yaml:"juhe_key"`
	PublicKey  string   `json:"-" yaml:"public_key"`
	PrivateKey string   `json:"-" yaml:"private_key"`
	AdminEmail []string `json:"admin_email"`
}
