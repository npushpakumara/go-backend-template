package config

import (
	"log"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
)

// Config represents the configuration for the application
// It includes settings for the server, database, JWT, logging, and AWS services.
type Config struct {
	Server  ServerConfig  `json:"server"`
	DB      DBConfig      `json:"db"`
	JWT     JWTConfig     `json:"jwt"`
	Logging LoggingConfig `json:"logging"`
	AWS     AWSConfigs    `json:"aws"`
}

// ServerConfig represents the configuration for the server
type ServerConfig struct {
	Port             uint          `json:"port"`
	Production       bool          `json:"production"`
	ReadTimeout      time.Duration `json:"read_timeout"`
	WriteTimeout     time.Duration `json:"write_timeout"`
	GracefulShutdown time.Duration `json:"graceful_shutdown"`
	Domain           string        `json:"domain"`
}

// DBConfig represents the configuration for the database
type DBConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
	SSLMode  string `json:"ssl_mode"`
	LogLevel int    `json:"log_level"`
	Pool     struct {
		MaxOpen     int           `json:"max_open"`
		MaxIdle     int           `json:"max_idle"`
		MaxLifetime time.Duration `json:"max_lifetime"`
	} `json:"pool"`
	MigrationEnabled bool `json:"migration_enabled"`
}

// JWTConfig represents the configuration for the JWT
type JWTConfig struct {
	Secret             string        `json:"secret"`
	RefreshTokenExpiry time.Duration `json:"refresh_token_exp"`
	AccessTokenExpiry  time.Duration `json:"access_token_exp"`
}

// LoggingConfig represents the configuration for logging
type LoggingConfig struct {
	Level    int    `json:"level"`
	Encoding string `json:"encoding"`
}

// AWSConfigs represents the configuration for AWS services
type AWSConfigs struct {
	Region    string `json:"region"`
	SESConfig struct {
		EmailFrom string `json:"from_email"`
	} `json:"ses"`
}

var k = koanf.New(".")

// LoadConfig loads the application configuration from environment variables and default settings.
// It initializes configuration using a default set of values and overrides them with environment variables.
func LoadConfig() (*Config, error) {
	// Load environment variables from a .env file into the application
	err := godotenv.Load()
	if err != nil {
		log.Printf("Cannot load .env file. err: %v, continue using system variables", err)
	}

	// Load default configuration settings
	err = k.Load(confmap.Provider(defaultConfigs, "."), nil)
	if err != nil {
		log.Printf("failed to load default config. err: %v", err)
		return nil, err
	}

	// Load environment variables with custom transformation
	transformKey := func(s string) string {
		n := 1
		if len(strings.Split(s, "_")) > 2 {
			n = 2
		}
		return strings.Replace(strings.ToLower(strings.TrimPrefix(s, "MYAPP_")), "_", ".", n)
	}

	// Load environment variables and apply custom transformation to keys
	if err := k.Load(env.Provider("MYAPP_", ".", transformKey), nil); err != nil {
		log.Printf("Failed to load config from environment variables: %v", err)
		return nil, err
	}

	// Unmarshal the loaded configuration into the Config struct
	var cfg Config
	if err := k.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{Tag: "json", FlatPaths: false}); err != nil {
		log.Printf("failed to unmarshal with conf. err: %v", err)
		return nil, err
	}

	return &cfg, err
}