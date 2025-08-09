package globals

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

// JWTKeyManager handles JWT RSA key pair management
type JWTKeyManager struct {
	config     *JWTConfig
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewJWTKeyManager creates a new JWT key manager
func NewJWTKeyManager(config *JWTConfig) *JWTKeyManager {
	return &JWTKeyManager{
		config: config,
	}
}

// InitializeKeys initializes JWT keys (load existing or generate new ones)
func (m *JWTKeyManager) InitializeKeys() error {
	err := m.ensureKeysDirectory()
	if err != nil {
		return fmt.Errorf("failed to create keys directory: %w", err)
	}

	privateKeyPath := m.getPrivateKeyPath()
	publicKeyPath := m.getPublicKeyPath()

	if m.keysExist() {
		slog.Info("Loading existing JWT keys", "private", privateKeyPath, "public", publicKeyPath)
		return m.loadExistingKeys()
	}

	if m.config.AutoGenerate {
		slog.Info("Generating new JWT keys", "key_size", m.config.RSAKeySize)
		return m.generateAndSaveKeys()
	}

	return fmt.Errorf("JWT keys not found and auto-generation is disabled")
}

// GenerateKeyPair generates a new RSA key pair
func (m *JWTKeyManager) GenerateKeyPair() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, m.config.RSAKeySize)
	if err != nil {
		return fmt.Errorf("failed to generate RSA private key: %w", err)
	}

	m.privateKey = privateKey
	m.publicKey = &privateKey.PublicKey

	return nil
}

// SaveKeys saves the generated keys to PEM files
func (m *JWTKeyManager) SaveKeys() error {
	if m.privateKey == nil {
		return fmt.Errorf("no private key to save")
	}

	err := m.savePrivateKey()
	if err != nil {
		return fmt.Errorf("failed to save private key: %w", err)
	}

	err = m.savePublicKey()
	if err != nil {
		return fmt.Errorf("failed to save public key: %w", err)
	}

	slog.Info("JWT keys saved successfully",
		"private_key", m.getPrivateKeyPath(),
		"public_key", m.getPublicKeyPath())

	return nil
}

// LoadKeys loads existing keys from PEM files
func (m *JWTKeyManager) LoadKeys() error {
	err := m.loadPrivateKey()
	if err != nil {
		return fmt.Errorf("failed to load private key: %w", err)
	}

	err = m.loadPublicKey()
	if err != nil {
		return fmt.Errorf("failed to load public key: %w", err)
	}

	slog.Info("JWT keys loaded successfully")
	return nil
}

// GetPrivateKeyPEM returns the private key in PEM format
func (m *JWTKeyManager) GetPrivateKeyPEM() ([]byte, error) {
	if m.privateKey == nil {
		return nil, fmt.Errorf("private key not initialized")
	}

	keyBytes, err := x509.MarshalPKCS8PrivateKey(m.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key: %w", err)
	}

	return pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	}), nil
}

// GetPublicKeyPEM returns the public key in PEM format
func (m *JWTKeyManager) GetPublicKeyPEM() ([]byte, error) {
	if m.publicKey == nil {
		return nil, fmt.Errorf("public key not initialized")
	}

	keyBytes, err := x509.MarshalPKIXPublicKey(m.publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}

	return pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: keyBytes,
	}), nil
}

// GetPrivateKey returns the RSA private key
func (m *JWTKeyManager) GetPrivateKey() *rsa.PrivateKey {
	return m.privateKey
}

// GetPublicKey returns the RSA public key
func (m *JWTKeyManager) GetPublicKey() *rsa.PublicKey {
	return m.publicKey
}

// RegenerateKeys generates new keys and saves them
func (m *JWTKeyManager) RegenerateKeys() error {
	slog.Info("Regenerating JWT keys")

	err := m.GenerateKeyPair()
	if err != nil {
		return fmt.Errorf("failed to generate new key pair: %w", err)
	}

	err = m.SaveKeys()
	if err != nil {
		return fmt.Errorf("failed to save new keys: %w", err)
	}

	return nil
}

// Private helper methods

func (m *JWTKeyManager) ensureKeysDirectory() error {
	return os.MkdirAll(m.config.KeysDirectory, 0700)
}

func (m *JWTKeyManager) getPrivateKeyPath() string {
	return filepath.Join(m.config.KeysDirectory, m.config.PrivateKeyFile)
}

func (m *JWTKeyManager) getPublicKeyPath() string {
	return filepath.Join(m.config.KeysDirectory, m.config.PublicKeyFile)
}

func (m *JWTKeyManager) keysExist() bool {
	privateKeyPath := m.getPrivateKeyPath()
	publicKeyPath := m.getPublicKeyPath()

	_, err1 := os.Stat(privateKeyPath)
	_, err2 := os.Stat(publicKeyPath)

	return err1 == nil && err2 == nil
}

func (m *JWTKeyManager) loadExistingKeys() error {
	err := m.loadPrivateKey()
	if err != nil {
		return err
	}

	err = m.loadPublicKey()
	if err != nil {
		return err
	}

	return nil
}

func (m *JWTKeyManager) generateAndSaveKeys() error {
	err := m.GenerateKeyPair()
	if err != nil {
		return err
	}

	return m.SaveKeys()
}

func (m *JWTKeyManager) savePrivateKey() error {
	keyPEM, err := m.GetPrivateKeyPEM()
	if err != nil {
		return err
	}

	return os.WriteFile(m.getPrivateKeyPath(), keyPEM, 0600)
}

func (m *JWTKeyManager) savePublicKey() error {
	keyPEM, err := m.GetPublicKeyPEM()
	if err != nil {
		return err
	}

	return os.WriteFile(m.getPublicKeyPath(), keyPEM, 0644)
}

func (m *JWTKeyManager) loadPrivateKey() error {
	keyData, err := os.ReadFile(m.getPrivateKeyPath())
	if err != nil {
		return fmt.Errorf("failed to read private key file: %w", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return fmt.Errorf("failed to decode PEM block for private key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS1 format as fallback
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse private key: %w", err)
		}
	}

	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return fmt.Errorf("private key is not RSA")
	}

	m.privateKey = rsaKey
	return nil
}

func (m *JWTKeyManager) loadPublicKey() error {
	keyData, err := os.ReadFile(m.getPublicKeyPath())
	if err != nil {
		return fmt.Errorf("failed to read public key file: %w", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return fmt.Errorf("failed to decode PEM block for public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("public key is not RSA")
	}

	m.publicKey = rsaKey
	return nil
}

// Global JWT key manager instance
var globalJWTKeyManager *JWTKeyManager

// InitializeJWTKeys initializes the global JWT key manager
func InitializeJWTKeys(config *JWTConfig) error {
	globalJWTKeyManager = NewJWTKeyManager(config)
	return globalJWTKeyManager.InitializeKeys()
}

// GetJWTKeyManager returns the global JWT key manager
func GetJWTKeyManager() *JWTKeyManager {
	return globalJWTKeyManager
}

// GetJWTPrivateKeyPEM returns the private key in PEM format
func GetJWTPrivateKeyPEM() ([]byte, error) {
	if globalJWTKeyManager == nil {
		return nil, fmt.Errorf("JWT key manager not initialized")
	}
	return globalJWTKeyManager.GetPrivateKeyPEM()
}

// GetJWTPublicKeyPEM returns the public key in PEM format
func GetJWTPublicKeyPEM() ([]byte, error) {
	if globalJWTKeyManager == nil {
		return nil, fmt.Errorf("JWT key manager not initialized")
	}
	return globalJWTKeyManager.GetPublicKeyPEM()
}
