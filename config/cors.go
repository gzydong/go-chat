package config

type Cors struct {
	Origin      string `json:"origin" yaml:"origin"`
	Headers     string `json:"headers" yaml:"headers"`
	Methods     string `json:"methods" yaml:"methods"`
	Credentials string `json:"credentials" yaml:"credentials"`
	MaxAge      string `json:"max_age" yaml:"max_age"`
}
