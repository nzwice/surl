package config

import (
	"time"

	"github.com/spf13/viper"
)

type DBConfig struct {
	DSN             string        `mapstructure:"dsn"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
}

type AppConfig struct {
	Debug    bool     `mapstructure:"debug"`
	DB       DBConfig `mapstructure:"db"`
	HttpAddr string   `mapstructure:"http_addr"`
	GrpcAddr string   `mapstructure:"grpc_addr"`
}

func Load(path string) (*AppConfig, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(path)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg AppConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
