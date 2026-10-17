package config

// Provider - абстракция над работой настроек для сервера (они могут быть получены из: ENV переменных, INI файла, БД, и т.п.)
//
//go:generate mockery
type Provider interface {
	Address() string
	AccrualSystemAddress() string
	JWTSecret() string
}

// Config - настройки для сервера
type Config struct {
	addr                 string
	accrualSystemAddress string
	jwtSecret            string
}

func New(cfg Provider) *Config {
	conf := &Config{}
	conf.AddressSet(cfg.Address())
	conf.AccrualSystemAddressSet(cfg.AccrualSystemAddress())
	conf.JWTSecretSet(cfg.JWTSecret())

	return conf
}

// Address - отвечает за адрес эндпоинта HTTP-сервера.
func (c *Config) Address() string {
	return c.addr
}
func (c *Config) AddressSet(addr string) {
	c.addr = addr
}

// AccrualSystemAddress - адрес системы расчёта начислений
func (c *Config) AccrualSystemAddress() string {
	return c.accrualSystemAddress
}
func (c *Config) AccrualSystemAddressSet(accrualSystemAddress string) {
	c.accrualSystemAddress = accrualSystemAddress
}

func (c *Config) JWTSecret() string {
	return c.jwtSecret
}

func (c *Config) JWTSecretSet(jwtSecret string) {
	c.jwtSecret = jwtSecret
}
