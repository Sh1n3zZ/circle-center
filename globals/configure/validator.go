package globals

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func ValidateConfig(config *Config) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	if err := validateServerConfig(&config.Server); err != nil {
		return fmt.Errorf("server config validation failed: %w", err)
	}

	if err := validateMySQLConfig(&config.MySQL); err != nil {
		return fmt.Errorf("mysql config validation failed: %w", err)
	}

	if err := validateRedisConfig(&config.Redis); err != nil {
		return fmt.Errorf("redis config validation failed: %w", err)
	}

	return nil
}

func validateServerConfig(config *ServerConfig) error {
	if config == nil {
		return fmt.Errorf("server config cannot be nil")
	}

	if config.Port <= 0 || config.Port > 65535 {
		return fmt.Errorf("invalid port number: %d (must be between 1 and 65535)", config.Port)
	}

	if config.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	if config.ReadTimeout <= 0 {
		return fmt.Errorf("read_timeout must be positive")
	}
	if config.WriteTimeout <= 0 {
		return fmt.Errorf("write_timeout must be positive")
	}
	if config.IdleTimeout <= 0 {
		return fmt.Errorf("idle_timeout must be positive")
	}

	if config.IdleTimeout <= config.ReadTimeout || config.IdleTimeout <= config.WriteTimeout {
		return fmt.Errorf("idle_timeout must be greater than read_timeout and write_timeout")
	}

	return nil
}

func validateMySQLConfig(config *MySQLConfig) error {
	if config == nil {
		return fmt.Errorf("mysql config cannot be nil")
	}

	if config.Host == "" {
		return fmt.Errorf("mysql host cannot be empty")
	}

	if config.Port <= 0 || config.Port > 65535 {
		return fmt.Errorf("invalid mysql port number: %d (must be between 1 and 65535)", config.Port)
	}

	if config.Username == "" {
		return fmt.Errorf("mysql username cannot be empty")
	}

	if config.Database == "" {
		return fmt.Errorf("mysql database name cannot be empty")
	}

	if config.Charset == "" {
		return fmt.Errorf("mysql charset cannot be empty")
	}

	if config.MaxOpenConns <= 0 {
		return fmt.Errorf("mysql max_open_conns must be positive")
	}
	if config.MaxIdleConns <= 0 {
		return fmt.Errorf("mysql max_idle_conns must be positive")
	}
	if config.MaxIdleConns > config.MaxOpenConns {
		return fmt.Errorf("mysql max_idle_conns cannot be greater than max_open_conns")
	}
	if config.MaxLifetime <= 0 {
		return fmt.Errorf("mysql max_lifetime must be positive")
	}

	if config.Loc == "" {
		return fmt.Errorf("mysql loc cannot be empty")
	}

	return nil
}

func validateRedisConfig(config *RedisConfig) error {
	if config == nil {
		return fmt.Errorf("redis config cannot be nil")
	}

	if config.Host == "" {
		return fmt.Errorf("redis host cannot be empty")
	}

	if config.Port <= 0 || config.Port > 65535 {
		return fmt.Errorf("invalid redis port number: %d (must be between 1 and 65535)", config.Port)
	}

	if config.DB < 0 || config.DB > 15 {
		return fmt.Errorf("invalid redis database number: %d (must be between 0 and 15)", config.DB)
	}

	if config.PoolSize <= 0 {
		return fmt.Errorf("redis pool_size must be positive")
	}
	if config.MinIdleConns < 0 {
		return fmt.Errorf("redis min_idle_conns cannot be negative")
	}
	if config.MinIdleConns > config.PoolSize {
		return fmt.Errorf("redis min_idle_conns cannot be greater than pool_size")
	}

	if config.MaxRetries < 0 {
		return fmt.Errorf("redis max_retries cannot be negative")
	}

	if config.DialTimeout <= 0 {
		return fmt.Errorf("redis dial_timeout must be positive")
	}
	if config.ReadTimeout <= 0 {
		return fmt.Errorf("redis read_timeout must be positive")
	}
	if config.WriteTimeout <= 0 {
		return fmt.Errorf("redis write_timeout must be positive")
	}
	if config.IdleTimeout <= 0 {
		return fmt.Errorf("redis idle_timeout must be positive")
	}

	return nil
}

// ValidateNetworkAddress validates if a host:port combination is valid
func ValidateNetworkAddress(host string, port int) error {
	if host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	if port <= 0 || port > 65535 {
		return fmt.Errorf("invalid port number: %d (must be between 1 and 65535)", port)
	}

	address := net.JoinHostPort(host, strconv.Itoa(port))
	if _, err := net.ResolveTCPAddr("tcp", address); err != nil {
		return fmt.Errorf("invalid network address %s: %w", address, err)
	}

	return nil
}

// ValidateDuration validates if a duration is positive
func ValidateDuration(d time.Duration, name string) error {
	if d <= 0 {
		return fmt.Errorf("%s must be positive", name)
	}
	return nil
}

// ValidatePositiveInt validates if an integer is positive
func ValidatePositiveInt(value int, name string) error {
	if value <= 0 {
		return fmt.Errorf("%s must be positive", name)
	}
	return nil
}

// ValidateNonNegativeInt validates if an integer is non-negative
func ValidateNonNegativeInt(value int, name string) error {
	if value < 0 {
		return fmt.Errorf("%s cannot be negative", name)
	}
	return nil
}

// ValidateRange validates if an integer is within a specified range
func ValidateRange(value, min, max int, name string) error {
	if value < min || value > max {
		return fmt.Errorf("%s must be between %d and %d", name, min, max)
	}
	return nil
}
