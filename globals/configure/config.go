package globals

import (
	"time"
)

// Config represents the main application configuration
type Config struct {
	Server ServerConfig `yaml:"server"`
	MySQL  MySQLConfig  `yaml:"mysql"`
	Redis  RedisConfig  `yaml:"redis"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         int           `yaml:"port"`
	Host         string        `yaml:"host"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

// MySQLConfig holds MySQL connection configuration
type MySQLConfig struct {
	Host           string        `yaml:"host"`
	Port           int           `yaml:"port"`
	Username       string        `yaml:"username"`
	Password       string        `yaml:"password"`
	Database       string        `yaml:"database"`
	Charset        string        `yaml:"charset"`
	ParseTime      bool          `yaml:"parse_time"`
	Loc            string        `yaml:"loc"`
	MaxOpenConns   int           `yaml:"max_open_conns"`
	MaxIdleConns   int           `yaml:"max_idle_conns"`
	MaxLifetime    time.Duration `yaml:"max_lifetime"`
	MultiStatement bool          `yaml:"multi_statement"`
}

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	Password     string        `yaml:"password"`
	DB           int           `yaml:"db"`
	PoolSize     int           `yaml:"pool_size"`
	MinIdleConns int           `yaml:"min_idle_conns"`
	MaxRetries   int           `yaml:"max_retries"`
	DialTimeout  time.Duration `yaml:"dial_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         8080,
			Host:         "0.0.0.0",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		MySQL: MySQLConfig{
			Host:           "localhost",
			Port:           3306,
			Username:       "root",
			Password:       "",
			Database:       "circle_center",
			Charset:        "utf8mb4",
			ParseTime:      true,
			Loc:            "Local",
			MaxOpenConns:   25,
			MaxIdleConns:   5,
			MaxLifetime:    time.Hour,
			MultiStatement: true,
		},
		Redis: RedisConfig{
			Host:         "localhost",
			Port:         6379,
			Password:     "",
			DB:           0,
			PoolSize:     10,
			MinIdleConns: 5,
			MaxRetries:   3,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			IdleTimeout:  5 * time.Minute,
		},
	}
}
