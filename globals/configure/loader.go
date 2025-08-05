package globals

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"sigs.k8s.io/yaml"
)

var (
	GlobalConfig *Config
)

func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = "config/config.yaml"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		slog.Warn("Failed to read config file, using default configuration", "path", configPath, "error", err)
		config := DefaultConfig()
		GlobalConfig = config
		return config, nil
	}

	config := DefaultConfig()
	if err := yaml.Unmarshal(data, config); err != nil {
		slog.Warn("Failed to parse config file, using default configuration", "path", configPath, "error", err)
		GlobalConfig = config
		return config, nil
	}

	if err := ValidateConfig(config); err != nil {
		slog.Warn("Config validation failed, using default configuration", "error", err)
		config = DefaultConfig()
		GlobalConfig = config
		return config, nil
	}

	GlobalConfig = config
	return config, nil
}

func LoadConfigFromBytes(data []byte) (*Config, error) {
	config := DefaultConfig()
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config data: %w", err)
	}

	if err := ValidateConfig(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	GlobalConfig = config
	return config, nil
}

func GetConfig() *Config {
	return GlobalConfig
}

func SaveConfig(config *Config, filePath string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file %s: %w", filePath, err)
	}

	return nil
}

func CreateDefaultConfig(filePath string) error {
	config := DefaultConfig()
	return SaveConfig(config, filePath)
}
