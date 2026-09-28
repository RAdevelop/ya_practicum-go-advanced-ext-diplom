package config

import "github.com/caarlos0/env/v11"

type Env struct {
	cfg envCfg
}

type envCfg struct {
	Addr string `env:"RUN_ADDRESS"`
}

func NewEnv() (*Env, error) {
	return NewEnvWithOptions(nil)
}

// NewEnvWithOptions - Конструктор с опциями
func NewEnvWithOptions(opts *env.Options) (*Env, error) {
	var cfg envCfg
	var err error

	if opts != nil {
		err = env.ParseWithOptions(&cfg, *opts)
	} else {
		err = env.Parse(&cfg)
	}

	if err != nil {
		return nil, err
	}

	return &Env{
		cfg: cfg,
	}, nil
}

func (env *Env) Address() string {
	return env.cfg.Addr
}
