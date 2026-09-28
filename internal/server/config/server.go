package config

// Provider - абстракция над работой настроек для сервера (они могут быть получены из: ENV переменных, INI файла, БД, и т.п.)
//
//go:generate mockery
type Provider interface {
	Address() string
}

// Config - настройки для сервера
type Config struct {
	addr string
}

func New(cfg Provider) *Config {
	conf := &Config{}
	conf.AddressSet(cfg.Address())

	return conf
}

// Address - отвечает за адрес эндпоинта HTTP-сервера.
func (c *Config) Address() string {
	return c.addr
}
func (c *Config) AddressSet(addr string) {
	c.addr = addr
}
