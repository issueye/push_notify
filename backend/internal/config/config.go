package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	AI       AIConfig       `mapstructure:"ai"`
	Webhook  WebhookConfig  `mapstructure:"webhook"`
}

type AppConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
}

type DatabaseConfig struct {
	Driver  string `mapstructure:"driver"`
	Path    string `mapstructure:"path"`
	LogMode string `mapstructure:"log_mode"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type JWTConfig struct {
	Secret             string `mapstructure:"secret"`
	AccessTokenExpire  int    `mapstructure:"access_token_expire"`
	RefreshTokenExpire int    `mapstructure:"refresh_token_expire"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

type AIConfig struct {
	DefaultModelID int `mapstructure:"default_model_id"`
	Timeout        int `mapstructure:"timeout"`
}

type WebhookConfig struct {
	SigningKey string `mapstructure:"signing_key"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 设置默认值
	if cfg.App.Port == 0 {
		cfg.App.Port = 8080
	}
	if cfg.JWT.AccessTokenExpire == 0 {
		cfg.JWT.AccessTokenExpire = 86400
	}
	if cfg.JWT.RefreshTokenExpire == 0 {
		cfg.JWT.RefreshTokenExpire = 604800
	}
	if cfg.AI.Timeout == 0 {
		cfg.AI.Timeout = 60
	}

	return &cfg, nil
}

func (c *Config) GetDSN() string {
	return c.Database.Path
}

func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}
