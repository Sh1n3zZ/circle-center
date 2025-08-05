package globals

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	// Test server config
	if config.Server.Port != 8080 {
		t.Errorf("Expected server port 8080, got %d", config.Server.Port)
	}

	// Test MySQL config
	if config.MySQL.Host != "localhost" {
		t.Errorf("Expected MySQL host localhost, got %s", config.MySQL.Host)
	}

	// Test Redis config
	if config.Redis.Host != "localhost" {
		t.Errorf("Expected Redis host localhost, got %s", config.Redis.Host)
	}
}

func TestLoadConfigFromBytes(t *testing.T) {
	yamlData := []byte(`
server:
  port: 9090
  host: "127.0.0.1"
mysql:
  host: "db.example.com"
  port: 3307
redis:
  host: "redis.example.com"
  port: 6380
`)

	config, err := LoadConfigFromBytes(yamlData)
	if err != nil {
		t.Fatalf("Failed to load config from bytes: %v", err)
	}

	if config.Server.Port != 9090 {
		t.Errorf("Expected server port 9090, got %d", config.Server.Port)
	}

	if config.MySQL.Host != "db.example.com" {
		t.Errorf("Expected MySQL host db.example.com, got %s", config.MySQL.Host)
	}

	if config.Redis.Host != "redis.example.com" {
		t.Errorf("Expected Redis host redis.example.com, got %s", config.Redis.Host)
	}
}

func TestValidateConfig(t *testing.T) {
	// Test valid config
	validConfig := DefaultConfig()
	if err := ValidateConfig(validConfig); err != nil {
		t.Errorf("Valid config should not have validation errors: %v", err)
	}

	// Test invalid server port
	invalidConfig := DefaultConfig()
	invalidConfig.Server.Port = 70000
	if err := ValidateConfig(invalidConfig); err == nil {
		t.Error("Invalid port should cause validation error")
	}

	// Test invalid MySQL host
	invalidConfig = DefaultConfig()
	invalidConfig.MySQL.Host = ""
	if err := ValidateConfig(invalidConfig); err == nil {
		t.Error("Empty MySQL host should cause validation error")
	}

	// Test invalid Redis port
	invalidConfig = DefaultConfig()
	invalidConfig.Redis.Port = 70000
	if err := ValidateConfig(invalidConfig); err == nil {
		t.Error("Invalid Redis port should cause validation error")
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	// Create a temporary file
	tmpFile := "test_config.yaml"
	defer os.Remove(tmpFile)

	// Create config and save it
	config := DefaultConfig()
	config.Server.Port = 9090
	config.MySQL.Host = "test-db"
	config.Redis.Host = "test-redis"

	if err := SaveConfig(config, tmpFile); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load the config back
	loadedConfig, err := LoadConfig(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify the loaded config matches the original
	if loadedConfig.Server.Port != config.Server.Port {
		t.Errorf("Server port mismatch: expected %d, got %d", config.Server.Port, loadedConfig.Server.Port)
	}

	if loadedConfig.MySQL.Host != config.MySQL.Host {
		t.Errorf("MySQL host mismatch: expected %s, got %s", config.MySQL.Host, loadedConfig.MySQL.Host)
	}

	if loadedConfig.Redis.Host != config.Redis.Host {
		t.Errorf("Redis host mismatch: expected %s, got %s", config.Redis.Host, loadedConfig.Redis.Host)
	}
}

func TestCreateDefaultConfig(t *testing.T) {
	// Create a temporary file
	tmpFile := "test_default_config.yaml"
	defer os.Remove(tmpFile)

	// Create default config file
	if err := CreateDefaultConfig(tmpFile); err != nil {
		t.Fatalf("Failed to create default config: %v", err)
	}

	// Verify the file was created
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Load and verify the config
	config, err := LoadConfig(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load created config: %v", err)
	}

	// Verify it has default values
	defaultConfig := DefaultConfig()
	if config.Server.Port != defaultConfig.Server.Port {
		t.Errorf("Expected default server port %d, got %d", defaultConfig.Server.Port, config.Server.Port)
	}
}
