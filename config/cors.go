package config

// Cors 跨域配置
type Cors struct {
	Origin      string `json:"origin" yaml:"origin"`
	Headers     string `json:"headers" yaml:"headers"`
	Methods     string `json:"methods" yaml:"methods"`
	Credentials string `json:"credentials" yaml:"credentials"`
	MaxAge      string `json:"max_age" yaml:"max_age"`
}

func (c *Cors) GetOrigin() string {
	return c.Origin
}

func (c *Cors) GetHeaders() string {
	return c.Headers
}

func (c *Cors) GetMethods() string {
	return c.Methods
}

func (c *Cors) GetCredentials() string {
	return c.Credentials
}

func (c *Cors) GetMaxAge() string {
	return c.MaxAge
}
