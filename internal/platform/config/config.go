package config

import "time"

type Config struct {
	Runtime Runtime `mapstructure:"runtime"`
	HTTP    HTTP    `mapstructure:"http"`
	DB      DB      `mapstructure:"db"`
	JWT     JWT     `mapstructure:"jwt"`
	Log     Log     `mapstructure:"log"`
}

type Runtime struct {
	Env string `mapstructure:"env"` // local, staging, production
}

type HTTP struct {
	Port              int           `mapstructure:"port"`
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
	ReadTimeout       time.Duration `mapstructure:"read_timeout"`
	WriteTimeout      time.Duration `mapstructure:"write_timeout"`
	IdleTimeout       time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdown_timeout"`
	MaxHeaderBytes    int           `mapstructure:"max_header_bytes"`
}

type DB struct {
	URL             string        `mapstructure:"url"`
	MaxConns        int32         `mapstructure:"max_conns"`
	MinConns        int32         `mapstructure:"min_conns"`
	MaxConnLifetime time.Duration `mapstructure:"max_conn_lifetime"`
	MaxConnIdleTime time.Duration `mapstructure:"max_conn_idle_time"`
}

type JWT struct {
	Secret string        `mapstructure:"secret"`
	TTL    time.Duration `mapstructure:"ttl"`
	Issuer string        `mapstructure:"issuer"`
}

type Log struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}
