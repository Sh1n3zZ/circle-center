-- Create users table migration
-- Core required fields, security enhancements, profile fields, and advanced features

CREATE TABLE users (
  -- Core required fields (all account systems need)
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  username VARCHAR(64) NOT NULL,
  email VARCHAR(255) NOT NULL,
  password_hash CHAR(100) NOT NULL COMMENT 'Password hash (bcrypt/argon2)',
  status TINYINT(1) UNSIGNED NOT NULL DEFAULT 0 COMMENT '0=inactive, 1=active, 2=disabled, 3=locked, 4=deleted',
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  
  -- Security enhancement fields (recommended to add)
  last_login_at TIMESTAMP(6) NULL,
  failed_attempts TINYINT UNSIGNED NOT NULL DEFAULT 0,
  locked_until TIMESTAMP(6) NULL,
  email_verified_at TIMESTAMP(6) NULL,
  password_changed_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  
  -- Profile fields (add based on business requirements)
  display_name VARCHAR(100) NULL,
  avatar_url VARCHAR(512) NULL,
  phone VARCHAR(20) NULL COMMENT 'E.164 format (+country code)',
  locale CHAR(5) NOT NULL DEFAULT 'en_US',
  timezone VARCHAR(50) NOT NULL DEFAULT 'UTC',
  
  -- Advanced feature fields (optional extensions)
  mfa_enabled BOOLEAN NOT NULL DEFAULT FALSE,
  mfa_secret VARBINARY(128) NULL COMMENT 'AES encrypted storage',
  recovery_codes TEXT NULL COMMENT 'Encrypted JSON array',
  privacy_version SMALLINT UNSIGNED NOT NULL DEFAULT 1,
  marketing_consent BOOLEAN NOT NULL DEFAULT FALSE,
  
  -- Indexes
  PRIMARY KEY (id),
  UNIQUE INDEX idx_unique_username (username),
  UNIQUE INDEX idx_unique_email (email),
  INDEX idx_status (status),
  INDEX idx_email_verified (email_verified_at),
  INDEX idx_phone (phone),
  INDEX idx_created_at (created_at),
  INDEX idx_last_login (last_login_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci 
  COMMENT='User account management table with security features and profile data';
