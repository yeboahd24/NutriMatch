package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// AppConfig represents the application configuration
type AppConfig struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DBConfig       `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Logging  LogConfig      `mapstructure:"logging"`
	Security SecurityConfig `mapstructure:"security"`
}

// ServerConfig represents the server configuration
type ServerConfig struct {
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

// DBConfig represents the database configuration
type DBConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"dbname"`
	SSLMode         string        `mapstructure:"sslmode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// JWTConfig represents the JWT configuration
type JWTConfig struct {
	AccessTokenSecret  string        `mapstructure:"access_token_secret"`
	RefreshTokenSecret string        `mapstructure:"refresh_token_secret"`
	AccessTokenExpiry  time.Duration `mapstructure:"access_token_expiry"`
	RefreshTokenExpiry time.Duration `mapstructure:"refresh_token_expiry"`
	Issuer             string        `mapstructure:"issuer"`
	Audience           string        `mapstructure:"audience"`
}

// LogConfig represents the logging configuration
type LogConfig struct {
	Level       string `mapstructure:"level"`
	Format      string `mapstructure:"format"`
	Output      string `mapstructure:"output"`
	ErrorOutput string `mapstructure:"error_output"`
}

// SecurityConfig represents security-related configuration
type SecurityConfig struct {
	ArgonMemory      uint32        `mapstructure:"argon_memory"`
	ArgonIterations  uint32        `mapstructure:"argon_iterations"`
	ArgonParallelism uint8         `mapstructure:"argon_parallelism"`
	ArgonSaltLength  uint32        `mapstructure:"argon_salt_length"`
	ArgonKeyLength   uint32        `mapstructure:"argon_key_length"`
	RateLimit        int           `mapstructure:"rate_limit"`
	RateLimitWindow  time.Duration `mapstructure:"rate_limit_window"`
}

// Load loads the configuration from files and environment variables
func Load() (*AppConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// Set default values
	setDefaults()

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found, will use defaults and env vars

	}

	// Override with environment variables
	viper.SetEnvPrefix("NUTRIMATCH")
	viper.AutomaticEnv()

	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return &config, nil
}

// setDefaults sets default values for configuration
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "5s")
	viper.SetDefault("server.write_timeout", "10s")
	viper.SetDefault("server.idle_timeout", "120s")
	viper.SetDefault("server.shutdown_timeout", "30s")

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.dbname", "nutrimatch")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 25)
	viper.SetDefault("database.conn_max_lifetime", "15m")

	// JWT defaults
	viper.SetDefault("jwt.access_token_expiry", "15m")
	viper.SetDefault("jwt.refresh_token_expiry", "7d")
	viper.SetDefault("jwt.issuer", "nutrimatch-api")
	viper.SetDefault("jwt.audience", "nutrimatch-clients")

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")
	viper.SetDefault("logging.error_output", "stderr")

	// Security defaults
	viper.SetDefault("security.argon_memory", 64*1024)
	viper.SetDefault("security.argon_iterations", 3)
	viper.SetDefault("security.argon_parallelism", 2)
	viper.SetDefault("security.argon_salt_length", 16)
	viper.SetDefault("security.argon_key_length", 32)
	viper.SetDefault("security.rate_limit", 100)
	viper.SetDefault("security.rate_limit_window", "1m")
}
