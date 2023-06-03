package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	RabbitMQ RabbitMQ `yaml:"rabbitmq"`
	MongoDB  MongoDB  `yaml:"mongodb"`
	Http     Http     `yaml:"http-server"`
	Logger   Logger   `yaml:"logger"`
}

type RabbitMQ struct {
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	Exchange    string `yaml:"exchange"`
	Queue       string `yaml:"queue"`
	RoutingKey  string `yaml:"routingKey"`
	ConsumerTag string `yaml:"consumerTag"`
	WorkerPool  int    `yaml:"workerPool"`
}

type MongoDB struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"db"`
}

type Http struct {
	AppVersion      string        `yaml:"appVersion"`
	Host            string        `yaml:"host"`
	Port            string        `yaml:"port"`
	ReadTimeout     time.Duration `yaml:"readTimeout"`
	WriteTimeout    time.Duration `yaml:"writeTimeout"`
	ShutdownTimeout time.Duration `yaml:"shutdownTimeout"`
}
type Logger struct {
	Development       bool   `yaml:"development"`
	DisableCaller     bool   `yaml:"disableCaller"`
	DisableStacktrace bool   `yaml:"disableStacktrace:"`
	Encoding          string `yaml:"encoding"`
	Level             string `yaml:"level"`
}

func NewConfig(cfgFile string) (*Config, error) {

	cfg := &Config{}

	if err := cleanenv.ReadConfig(cfgFile, cfg); err != nil {
		return nil, err
	}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
