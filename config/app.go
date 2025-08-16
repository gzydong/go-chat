package config

type App struct {
	Env        string   `yaml:"env"`
	Debug      bool     `yaml:"debug"`
	PublicKey  string   `yaml:"public_key"`
	PrivateKey string   `yaml:"private_key"`
	AesKey     string   `yaml:"aes_key"`
	AdminEmail []string `yaml:"admin_email"`
}
