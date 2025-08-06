# Database Repository with SQLC

This directory contains the database schema, migrations, and SQLC configuration for generating type-safe Go code for database operations.

## Structure

- `migrations/` - Database migration files
- `query/` - SQL query files for SQLC
- `sqlc/` - Generated Go code (created by SQLC)
- `sqlc.yaml` - SQLC configuration file

## Prerequisites

1. Install SQLC CLI:

   ```bash
   go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
   ```

2. Ensure you have MySQL database running and accessible

## Usage

### Generate Go Code

To generate type-safe Go code from SQL queries:

```bash
# From the project root directory
sqlc generate --file repository/sqlc.yaml
```

This command will:

- Read the SQL queries from `repository/query/`
- Read the database schema from `repository/migrations/v1_account/`
- Generate Go code in `repository/sqlc/account/`

### Verify SQLC Configuration

To validate your SQLC configuration:

```bash
sqlc validate
```

### Generate Code for Specific Database

If you need to generate code for a specific database engine:

```bash
sqlc generate --config repository/sqlc.yaml
```

## Generated Files

After running `sqlc generate`, the following files will be created in `repository/sqlc/account/`:

- `db.go` - Database connection and interface
- `models.go` - Generated structs for database tables
- `querier.go` - Interface for database operations
- `account.sql.go` - Generated methods for SQL queries

## Database Operations

The generated code provides the following operations for the `users` table:

### User Management

- `CreateUser` - Create a new user account
- `GetUserByID` - Retrieve user by ID
- `GetUserByUsername` - Retrieve user by username
- `GetUserByEmail` - Retrieve user by email
- `GetUserByPhone` - Retrieve user by phone number
- `ListUsers` - List all active users with pagination
- `UpdateUserProfile` - Update user profile information
- `UpdateUserPassword` - Update user password
- `DeleteUser` - Soft delete user (set status to deleted)

### Authentication & Security

- `UpdateLastLogin` - Update last login timestamp
- `IncrementFailedAttempts` - Increment failed login attempts
- `LockUser` - Lock user account
- `UnlockUser` - Unlock user account
- `VerifyEmail` - Mark email as verified

### Advanced Features

- `UpdateMFASettings` - Update MFA settings
- `UpdateMarketingConsent` - Update marketing consent
- `UpdatePrivacyVersion` - Update privacy policy version

### Queries & Analytics

- `CountUsers` - Count total active users
- `CountUsersByStatus` - Count users by status
- `GetUsersByEmailVerified` - Get users with verified emails
- `GetUsersByLastLogin` - Get users by last login time
- `GetLockedUsers` - Get currently locked users
- `GetUsersWithMFA` - Get users with MFA enabled

## Example Usage

```go
package main

import (
    "context"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    
    "your-project/repository/sqlc/account"
)

func main() {
    db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/dbname")
    if err != nil {
        panic(err)
    }
    defer db.Close()
    
    queries := account.New(db)
    ctx := context.Background()
    
    // Create a new user
    result, err := queries.CreateUser(ctx, account.CreateUserParams{
        Username:     "john_doe",
        Email:        "john@example.com",
        PasswordHash: "hashed_password",
        DisplayName:  "John Doe",
        Phone:        "+1234567890",
        Locale:       "en_US",
        Timezone:     "America/New_York",
    })
    
    // Get user by ID
    user, err := queries.GetUserByID(ctx, 1)
    
    // Update user profile
    err = queries.UpdateUserProfile(ctx, account.UpdateUserProfileParams{
        DisplayName: "John Smith",
        AvatarUrl:   "https://example.com/avatar.jpg",
        Phone:       "+1234567890",
        Locale:      "en_US",
        Timezone:    "America/New_York",
        ID:          1,
    })
}
```

## Migration Management

To apply database migrations, use your preferred migration tool (e.g., golang-migrate, Atlas, etc.) with the files in `migrations/`.

## Notes

- All queries automatically exclude deleted users (status = 4)
- Timestamps are automatically updated on modifications
- The generated code includes proper error handling and type safety
- Use prepared statements for better performance and security
