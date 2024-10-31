package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server" validate:"required"`
	Database DatabaseConfig `mapstructure:"database" validate:"required"`
	Import   ImportConfig   `mapstructure:"import" validate:"required"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

type ServerConfig struct {
	Port int `mapstructure:"port" validate:"required,min=1,max=65535"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host" validate:"required"`
	Port     int    `mapstructure:"port" validate:"required,min=1,max=65535"`
	Username string `mapstructure:"username" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
	DBName   string `mapstructure:"dbname" validate:"required"`
	SSLMode  string `mapstructure:"sslmode" validate:"required,oneof=disable enable verify-full"`
}

type JWTConfig struct {
	SecretKey   string `mapstructure:"secret_key" validate:"required"`
	ExpiryHours int    `mapstructure:"expiry_hours" validate:"required,min=1"`
}

type ImportConfig struct {
	Directory string `mapstructure:"directory" validate:"required,dir"`
}

func InitializeConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Password, c.DBName, c.SSLMode)
}
