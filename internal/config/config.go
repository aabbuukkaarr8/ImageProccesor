package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LogLevel string       `yaml:"loglevel"`
	Server   ServerConfig `yaml:"server"`
	DB       DBConfig     `yaml:"db"`
	MinIO    MinIOConfig  `yaml:"minio"`
	Kafka    KafkaConfig  `yaml:"kafka"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`

	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime string `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime string `yaml:"conn_max_idle_time"`
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode,
	)
}

func (c DBConfig) ParseConnMaxLifetime() time.Duration {
	if c.ConnMaxLifetime == "" {
		return 0
	}
	d, err := time.ParseDuration(c.ConnMaxLifetime)
	if err != nil {
		return 0
	}
	return d
}

func (c DBConfig) ParseConnMaxIdleTime() time.Duration {
	if c.ConnMaxIdleTime == "" {
		return 0
	}
	d, err := time.ParseDuration(c.ConnMaxIdleTime)
	if err != nil {
		return 0
	}
	return d
}

// MinIOConfig хранит конфигурацию MinIO
type MinIOConfig struct {
	Endpoint  string `yaml:"endpoint"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Bucket    string `yaml:"bucket"`
	UseSSL    bool   `yaml:"use_ssl"`
}

// KafkaConfig хранит конфигурацию Kafka
type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
	GroupID string   `yaml:"group_id"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
