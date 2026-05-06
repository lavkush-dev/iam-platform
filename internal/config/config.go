package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App    AppConfig    `mapstructure:"app"`
	Server ServerConfig `mapstructure:"server"`
	DB     DBConfig     `mapstructure:"db"`
	Redis  RedisConfig  `mapstructure:"redis"`
	JWT    JWTConfig    `mapstructure:"jwt"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type DBConfig struct {
	URL string `mapstructure:"url"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	Expiration int    `mapstructure:"expiration"` // minutes
}

func Load() *Config {
	v := viper.New()

	// Config file
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./configs")

	// Environment variables override
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Defaults (safe fallbacks)
	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		log.Println("no config file found, using env/defaults")
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("unable to decode config: %v", err)
	}

	return &cfg
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("server.port", "8080")

	v.SetDefault("db.url", "postgres://iam:secret@localhost:5432/iam_db?sslmode=disable")

	v.SetDefault("redis.addr", "localhost:6379")
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)

	v.SetDefault("jwt.secret", "supersecret")
	v.SetDefault("jwt.expiration", 60)
}
