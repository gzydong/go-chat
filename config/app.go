package config

type App struct {
	Env        string   `json:"env"`
	Debug      bool     `json:"debug"`
	PublicKey  string   `json:"-" yaml:"public_key"`
	PrivateKey string   `json:"-" yaml:"private_key"`
	AdminEmail []string `json:"admin_email"`
}
