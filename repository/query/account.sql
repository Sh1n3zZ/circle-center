-- name: CreateUser :execresult
INSERT INTO users (
  username, email, password_hash, display_name, phone, locale, timezone
) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetUserByID :one
SELECT * FROM users WHERE id = ? AND status != 4 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = ? AND status != 4 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ? AND status != 4 LIMIT 1;

-- name: GetUserByPhone :one
SELECT * FROM users WHERE phone = ? AND status != 4 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users WHERE status != 4 ORDER BY created_at DESC LIMIT ? OFFSET ?;

-- name: ListUsersByStatus :many
SELECT * FROM users WHERE status = ? ORDER BY created_at DESC LIMIT ? OFFSET ?;

-- name: UpdateUserProfile :exec
UPDATE users SET 
  display_name = ?, 
  avatar_url = ?, 
  phone = ?, 
  locale = ?, 
  timezone = ?,
  updated_at = CURRENT_TIMESTAMP(6)
WHERE id = ?;

-- name: UpdateUserPassword :exec
UPDATE users SET 
  password_hash = ?, 
  password_changed_at = CURRENT_TIMESTAMP(6),
  updated_at = CURRENT_TIMESTAMP(6)
WHERE id = ?;

-- name: UpdateUserStatus :exec
UPDATE users SET 
  status = ?, 
  updated_at = CURRENT_TIMESTAMP(6)
WHERE id = ?;

-- name: UpdateLastLogin :exec
UPDATE users SET 
  last_login_at = CURRENT_TIMESTAMP(6),
  failed_attempts = 0,
  locked_until = NULL,
  updated_at = CURRENT_TIMESTAMP(6)
WHERE id = ?;

-- name: IncrementFailedAttempts :exec
UPDATE users SET 
  failed_attempts = failed_attempts + 1,
  updated_at = CURRENT_TIMESTAMP(6)
WHERE id = ?;

-- name: LockUser :exec
UPDATE users SET 
  status = 3,
  locked_until = ?,
  updated_at = CURRENT_TIMESTAMP(6)
WHERE id = ?;

-- name: UnlockUser :exec
UPDATE users SET 
  status = 1,
  locked_until = NULL,
  failed_attempts = 0,
  updated_at = CURRENT_TIMESTAMP(6)
WHERE id = ?;

-- name: VerifyEmail :exec
UPDATE users SET 
  email_verified_at = CURRENT_TIMESTAMP(6),
  updated_at = CURRENT_TIMESTAMP(6)
WHERE id = ?;

-- name: UpdateMFASettings :exec
UPDATE users SET 
  mfa_enabled = ?,
  mfa_secret = ?,
  recovery_codes = ?,
  updated_at = CURRENT_TIMESTAMP(6)
WHERE id = ?;

-- name: UpdateMarketingConsent :exec
UPDATE users SET 
  marketing_consent = ?,
  updated_at = CURRENT_TIMESTAMP(6)
WHERE id = ?;

-- name: UpdatePrivacyVersion :exec
UPDATE users SET 
  privacy_version = ?,
  updated_at = CURRENT_TIMESTAMP(6)
WHERE id = ?;

-- name: DeleteUser :exec
UPDATE users SET 
  status = 4,
  updated_at = CURRENT_TIMESTAMP(6)
WHERE id = ?;

-- name: CountUsers :one
SELECT COUNT(*) FROM users WHERE status != 4;

-- name: CountUsersByStatus :one
SELECT COUNT(*) FROM users WHERE status = ?;

-- name: GetUsersByEmailVerified :many
SELECT * FROM users WHERE email_verified_at IS NOT NULL AND status != 4 ORDER BY created_at DESC;

-- name: GetUsersByLastLogin :many
SELECT * FROM users WHERE last_login_at IS NOT NULL AND status != 4 ORDER BY last_login_at DESC LIMIT ? OFFSET ?;

-- name: GetLockedUsers :many
SELECT * FROM users WHERE status = 3 AND locked_until IS NOT NULL ORDER BY locked_until ASC;

-- name: GetUsersWithMFA :many
SELECT * FROM users WHERE mfa_enabled = TRUE AND status != 4 ORDER BY created_at DESC;
