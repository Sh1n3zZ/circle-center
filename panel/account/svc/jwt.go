package account

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"

	configure "circle-center/globals/configure"
)

// JWTClient handles JWT token generation and validation
type JWTClient struct {
	privateKey jwk.Key
	publicKey  jwk.Key
	issuer     string
	expiryTime time.Duration
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Issuer     string
	ExpiryTime time.Duration
	RSAKeySize int
}

// DefaultJWTConfig returns default JWT configuration
func DefaultJWTConfig() *JWTConfig {
	return &JWTConfig{
		Issuer:     "circle-center",
		ExpiryTime: 24 * time.Hour,
		RSAKeySize: 2048,
	}
}

// NewJWTClient creates a new JWT client with RSA key pair
func NewJWTClient(config *JWTConfig) (*JWTClient, error) {
	if config == nil {
		config = DefaultJWTConfig()
	}

	// Generate RSA private key
	privateKeyRSA, err := rsa.GenerateKey(rand.Reader, config.RSAKeySize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}

	// Convert to JWK private key
	privateKey, err := jwk.Import(privateKeyRSA)
	if err != nil {
		return nil, fmt.Errorf("failed to create private JWK: %w", err)
	}

	// Set key ID and algorithm
	err = privateKey.Set(jwk.KeyIDKey, "default")
	if err != nil {
		return nil, fmt.Errorf("failed to set key ID: %w", err)
	}

	err = privateKey.Set(jwk.AlgorithmKey, jwa.RS256())
	if err != nil {
		return nil, fmt.Errorf("failed to set algorithm: %w", err)
	}

	// Get public key
	publicKey, err := jwk.PublicKeyOf(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to extract public key: %w", err)
	}

	return &JWTClient{
		privateKey: privateKey,
		publicKey:  publicKey,
		issuer:     config.Issuer,
		expiryTime: config.ExpiryTime,
	}, nil
}

// NewJWTClientFromGlobalKeys creates a new JWT client using global key manager
func NewJWTClientFromGlobalKeys() (*JWTClient, error) {
	keyManager := configure.GetJWTKeyManager()
	if keyManager == nil {
		return nil, fmt.Errorf("JWT key manager not initialized")
	}

	// Get keys from global key manager
	privateKeyPEM, err := keyManager.GetPrivateKeyPEM()
	if err != nil {
		return nil, fmt.Errorf("failed to get private key from key manager: %w", err)
	}

	publicKeyPEM, err := keyManager.GetPublicKeyPEM()
	if err != nil {
		return nil, fmt.Errorf("failed to get public key from key manager: %w", err)
	}

	// Get JWT config for issuer and expiry time
	config := configure.GetConfig()
	if config == nil {
		return nil, fmt.Errorf("configuration not loaded")
	}

	jwtConfig := &JWTConfig{
		Issuer:     config.JWT.Issuer,
		ExpiryTime: config.JWT.ExpiryTime,
		RSAKeySize: config.JWT.RSAKeySize,
	}

	return NewJWTClientWithKeys(privateKeyPEM, publicKeyPEM, jwtConfig)
}

// NewJWTClientWithKeys creates a new JWT client with existing keys
func NewJWTClientWithKeys(privateKeyPEM, publicKeyPEM []byte, config *JWTConfig) (*JWTClient, error) {
	if config == nil {
		config = DefaultJWTConfig()
	}

	// Parse private key
	privateKey, err := parsePrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Parse public key
	publicKey, err := parsePublicKeyFromPEM(publicKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return &JWTClient{
		privateKey: privateKey,
		publicKey:  publicKey,
		issuer:     config.Issuer,
		expiryTime: config.ExpiryTime,
	}, nil
}

// GenerateToken generates a JWT token for a user
func (j *JWTClient) GenerateToken(userID uint64, username, email string) (string, error) {
	now := time.Now()

	// Build JWT token
	token, err := jwt.NewBuilder().
		Issuer(j.issuer).
		Subject(fmt.Sprintf("%d", userID)).
		Audience([]string{"circle-center-users"}).
		IssuedAt(now).
		NotBefore(now).
		Expiration(now.Add(j.expiryTime)).
		Claim("user_id", userID).
		Claim("username", username).
		Claim("email", email).
		Build()
	if err != nil {
		return "", fmt.Errorf("failed to build token: %w", err)
	}

	// Sign the token
	signed, err := jwt.Sign(token, jwt.WithKey(jwa.RS256(), j.privateKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return string(signed), nil
}

// ValidateToken validates a JWT token and returns the claims
func (j *JWTClient) ValidateToken(tokenString string) (jwt.Token, error) {
	// Parse and verify the token
	token, err := jwt.Parse([]byte(tokenString), jwt.WithKey(jwa.RS256(), j.publicKey))
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	return token, nil
}

// ExtractUserInfo extracts user information from a validated token
func (j *JWTClient) ExtractUserInfo(token jwt.Token) (userID uint64, username, email string, err error) {
	// Extract user ID
	var userIDClaim interface{}
	if err := token.Get("user_id", &userIDClaim); err != nil {
		return 0, "", "", fmt.Errorf("user_id claim not found: %w", err)
	}

	// Convert to uint64
	switch v := userIDClaim.(type) {
	case float64:
		userID = uint64(v)
	case int64:
		userID = uint64(v)
	case uint64:
		userID = v
	default:
		return 0, "", "", fmt.Errorf("invalid user_id type")
	}

	// Extract username
	if err := token.Get("username", &username); err != nil {
		return 0, "", "", fmt.Errorf("username claim not found: %w", err)
	}

	// Extract email
	if err := token.Get("email", &email); err != nil {
		return 0, "", "", fmt.Errorf("email claim not found: %w", err)
	}

	return userID, username, email, nil
}

// GetPublicKey returns the public key in JWK format
func (j *JWTClient) GetPublicKey() jwk.Key {
	return j.publicKey
}

// GetPrivateKey returns the private key in JWK format
func (j *JWTClient) GetPrivateKey() jwk.Key {
	return j.privateKey
}

// ExportPrivateKeyPEM exports the private key as PEM format
func (j *JWTClient) ExportPrivateKeyPEM() ([]byte, error) {
	var rawKey interface{}
	err := jwk.Export(j.privateKey, &rawKey)
	if err != nil {
		return nil, fmt.Errorf("failed to export private key: %w", err)
	}

	rsaKey, ok := rawKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not RSA")
	}

	keyBytes, err := x509.MarshalPKCS8PrivateKey(rsaKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key: %w", err)
	}

	return pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	}), nil
}

// ExportPublicKeyPEM exports the public key as PEM format
func (j *JWTClient) ExportPublicKeyPEM() ([]byte, error) {
	var rawKey interface{}
	err := jwk.Export(j.publicKey, &rawKey)
	if err != nil {
		return nil, fmt.Errorf("failed to export public key: %w", err)
	}

	rsaKey, ok := rawKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not RSA")
	}

	keyBytes, err := x509.MarshalPKIXPublicKey(rsaKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}

	return pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: keyBytes,
	}), nil
}

// parsePrivateKeyFromPEM parses a private key from PEM format
func parsePrivateKeyFromPEM(keyPEM []byte) (jwk.Key, error) {
	block, _ := pem.Decode(keyPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS1 format
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
	}

	jwkKey, err := jwk.Import(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWK from private key: %w", err)
	}

	// Set key ID and algorithm
	err = jwkKey.Set(jwk.KeyIDKey, "default")
	if err != nil {
		return nil, fmt.Errorf("failed to set key ID: %w", err)
	}

	err = jwkKey.Set(jwk.AlgorithmKey, jwa.RS256())
	if err != nil {
		return nil, fmt.Errorf("failed to set algorithm: %w", err)
	}

	return jwkKey, nil
}

// parsePublicKeyFromPEM parses a public key from PEM format
func parsePublicKeyFromPEM(keyPEM []byte) (jwk.Key, error) {
	block, _ := pem.Decode(keyPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	jwkKey, err := jwk.Import(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWK from public key: %w", err)
	}

	// Set key ID and algorithm
	err = jwkKey.Set(jwk.KeyIDKey, "default")
	if err != nil {
		return nil, fmt.Errorf("failed to set key ID: %w", err)
	}

	err = jwkKey.Set(jwk.AlgorithmKey, jwa.RS256())
	if err != nil {
		return nil, fmt.Errorf("failed to set algorithm: %w", err)
	}

	return jwkKey, nil
}
