package globals

import (
	"time"
)

// Config represents the main application configuration
type Config struct {
	Server   ServerConfig   `yaml:"server" json:"server"`
	MySQL    MySQLConfig    `yaml:"mysql" json:"mysql"`
	Redis    RedisConfig    `yaml:"redis" json:"redis"`
	Mail     MailConfig     `yaml:"mail" json:"mail"`
	Frontend FrontendConfig `yaml:"frontend" json:"frontend"`
	JWT      JWTConfig      `yaml:"jwt" json:"jwt"`
	Avatar   AvatarConfig   `yaml:"avatar" json:"avatar"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         int           `yaml:"port" json:"port"`
	Host         string        `yaml:"host" json:"host"`
	ReadTimeout  time.Duration `yaml:"read_timeout" json:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout" json:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" json:"idle_timeout"`
}

// MySQLConfig holds MySQL connection configuration
type MySQLConfig struct {
	Host           string        `yaml:"host" json:"host"`
	Port           int           `yaml:"port" json:"port"`
	Username       string        `yaml:"username" json:"username"`
	Password       string        `yaml:"password" json:"password"`
	Database       string        `yaml:"database" json:"database"`
	Charset        string        `yaml:"charset" json:"charset"`
	ParseTime      bool          `yaml:"parse_time" json:"parse_time"`
	Loc            string        `yaml:"loc" json:"loc"`
	MaxOpenConns   int           `yaml:"max_open_conns" json:"max_open_conns"`
	MaxIdleConns   int           `yaml:"max_idle_conns" json:"max_idle_conns"`
	MaxLifetime    time.Duration `yaml:"max_lifetime" json:"max_lifetime"`
	MultiStatement bool          `yaml:"multi_statement" json:"multi_statement"`
}

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Host         string        `yaml:"host" json:"host"`
	Port         int           `yaml:"port" json:"port"`
	Password     string        `yaml:"password" json:"password"`
	DB           int           `yaml:"db" json:"db"`
	PoolSize     int           `yaml:"pool_size" json:"pool_size"`
	MinIdleConns int           `yaml:"min_idle_conns" json:"min_idle_conns"`
	MaxRetries   int           `yaml:"max_retries" json:"max_retries"`
	DialTimeout  time.Duration `yaml:"dial_timeout" json:"dial_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout" json:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout" json:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" json:"idle_timeout"`
}

// MailConfig holds email configuration
type MailConfig struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	From     string `yaml:"from" json:"from"`
	TLSMode  string `yaml:"tls_mode" json:"tls_mode"` // "mandatory" (STARTTLS), "opportunistic" (STARTTLS with fallback), "ssl" (SSL), "none" (NoTLS)
}

// FrontendConfig holds frontend configuration
type FrontendConfig struct {
	BaseURL string `yaml:"base_url" json:"base_url"`
}

// AvatarConfig holds avatar related configuration
type AvatarConfig struct {
	// Maximum allowed avatar upload size in bytes
	MaxUploadBytes int64 `yaml:"max_upload_bytes" json:"max_upload_bytes"`
	// Maximum image quality allowed when processing/returning images (1-100)
	MaxImageQuality int `yaml:"max_image_quality" json:"max_image_quality"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Issuer         string        `yaml:"issuer" json:"issuer"`
	ExpiryTime     time.Duration `yaml:"expiry_time" json:"expiry_time"`
	RSAKeySize     int           `yaml:"rsa_key_size" json:"rsa_key_size"`
	KeysDirectory  string        `yaml:"keys_directory" json:"keys_directory"`
	PrivateKeyFile string        `yaml:"private_key_file" json:"private_key_file"`
	PublicKeyFile  string        `yaml:"public_key_file" json:"public_key_file"`
	AutoGenerate   bool          `yaml:"auto_generate" json:"auto_generate"`
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
		Mail: MailConfig{
			Host:     "smtp.gmail.com",
			Port:     587,
			Username: "",
			Password: "",
			From:     "no_reply@gmail.com",
			TLSMode:  "ssl",
		},
		Frontend: FrontendConfig{
			BaseURL: "http://localhost:5173",
		},
		Avatar: AvatarConfig{
			MaxUploadBytes:  2 * 1024 * 1024, // 2MB
			MaxImageQuality: 90,
		},
		JWT: JWTConfig{
			Issuer:         "circle-center",
			ExpiryTime:     168 * time.Hour,
			RSAKeySize:     2048,
			KeysDirectory:  "config/keys",
			PrivateKeyFile: "jwt_private.pem",
			PublicKeyFile:  "jwt_public.pem",
			AutoGenerate:   true,
		},
	}
}
