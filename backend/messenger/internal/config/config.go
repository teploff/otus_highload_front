package config

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HttpAddr   string           `mapstructure:"http_addr"`
	GRPCAddr   string           `mapstructure:"grpc_addr"`
	Auth       AuthConfig       `mapstructure:"auth"`
	Clickhouse ClickhouseConfig `mapstructure:"ch"`
	Jaeger     JaegerConfig     `mapstructure:"jaeger"`
	Consul     ConsulConfig     `mapstructure:"consul"`
	Zabbix     ZabbixConfig     `mapstructure:"zabbix"`
	Logger     LoggerConfig     `mapstructure:"logger"`
	Sharding   ShardingConfig   `mapstructure:"sharding"`
}

// AuthConfig authorization service address.
//
// Addr - socket address (host + port).
type AuthConfig struct {
	Addr string `mapstructure:"addr"`
}

type ClickhouseConfig struct {
	DSN         string        `mapstructure:"dsn"`
	PushTimeout time.Duration `mapstructure:"push_timeout"`
}

// LoggerConfig logger configuration.
type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

type JaegerConfig struct {
	AgentAddr   string `mapstructure:"agent_addr"`
	ServiceName string `mapstructure:"service_name"`
}

type ConsulConfig struct {
	Addr        string `mapstructure:"addr"`
	AgentAddr   string `mapstructure:"agent_addr"`
	ServiceID   string `mapstructure:"service_id"`
	ServiceName string `mapstructure:"service_name"`
}

type ZabbixConfig struct {
	ServerHost string `mapstructure:"server_host"`
	Port       int    `mapstructure:"port"`
	HostName   string `mapstructure:"host_name"`
}

type ShardingConfig struct {
	CountNodes int `mapstructure:"count_nodes"`
}

// Load create configuration from file & environments.
func Load(path string) (*Config, error) {
	dir, file := filepath.Split(path)
	viper.SetConfigName(strings.TrimSuffix(file, filepath.Ext(file)))
	viper.AddConfigPath(dir)

	var cfg Config

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file, %w", err)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("fail to decode into struct, %w", err)
	}

	return &cfg, nil
}
