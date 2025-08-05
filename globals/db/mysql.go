package globals

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

var (
	DB *sqlx.DB
)

type MySQLConfig struct {
	Host           string
	Port           int
	Username       string
	Password       string
	Database       string
	Charset        string
	ParseTime      bool
	Loc            string
	MaxOpenConns   int
	MaxIdleConns   int
	MaxLifetime    time.Duration
	MultiStatement bool
}

// convertPathToURL converts Windows backslashes to forward slashes for URL format
func convertPathToURL(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

// ConnectMySQL establishes a connection to MySQL database
func ConnectMySQL(config *MySQLConfig) error {
	if config == nil {
		slog.Error("MySQL configuration is required")
		os.Exit(1)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.Charset,
		config.ParseTime,
		config.Loc,
	)

	if config.MultiStatement {
		dsn += "&multiStatements=true"
		slog.Info("Multi-statement support enabled in DSN")
	}

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.MaxLifetime)

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping MySQL: %w", err)
	}

	DB = db
	slog.Info("Successfully connected to MySQL database",
		"host", config.Host,
		"port", config.Port,
		"database", config.Database,
		"max_open_conns", config.MaxOpenConns,
		"max_idle_conns", config.MaxIdleConns,
	)
	return nil
}

// CloseMySQL closes the database connection
func CloseMySQL() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// GetDB returns the global database connection
func GetDB() *sqlx.DB {
	return DB
}

// RunMigrations runs database migrations from the specified path and its subdirectories
func RunMigrations(migrationsPath string) error {
	if DB == nil {
		return fmt.Errorf("database connection not established")
	}

	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		slog.Info("Migrations directory does not exist, skipping migrations",
			"migrations_path", migrationsPath,
		)
		return nil
	}

	sqlDB := DB.DB
	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("failed to create MySQL driver: %w", err)
	}

	entries, err := os.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var totalMigrations int
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		subDirPath := filepath.Join(migrationsPath, entry.Name())

		subEntries, err := os.ReadDir(subDirPath)
		if err != nil {
			slog.Warn("Failed to read subdirectory", "path", subDirPath, "error", err)
			continue
		}

		hasMigrationFiles := false
		for _, subEntry := range subEntries {
			if !subEntry.IsDir() && strings.HasSuffix(subEntry.Name(), ".up.sql") {
				hasMigrationFiles = true
				break
			}
		}

		if !hasMigrationFiles {
			continue
		}

		m, err := migrate.NewWithDatabaseInstance(
			fmt.Sprintf("file://%s", convertPathToURL(subDirPath)),
			"mysql",
			driver,
		)
		if err != nil {
			slog.Warn("Failed to create migrate instance for subdirectory",
				"path", subDirPath, "error", err)
			continue
		}

		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			slog.Warn("Failed to run migrations for subdirectory",
				"path", subDirPath, "error", err)
		} else {
			totalMigrations++
			slog.Info("Successfully ran migrations for subdirectory", "path", subDirPath)
		}
		m.Close()
	}

	rootEntries, err := os.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	hasRootMigrations := false
	for _, entry := range rootEntries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			hasRootMigrations = true
			break
		}
	}

	if hasRootMigrations {
		m, err := migrate.NewWithDatabaseInstance(
			fmt.Sprintf("file://%s", convertPathToURL(migrationsPath)),
			"mysql",
			driver,
		)
		if err != nil {
			return fmt.Errorf("failed to create migrate instance for root directory: %w", err)
		}
		defer m.Close()

		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("failed to run migrations: %w", err)
		}
		totalMigrations++
	}

	if totalMigrations == 0 {
		slog.Info("No migration files found",
			"migrations_path", migrationsPath,
		)
	} else {
		slog.Info("Database migrations completed successfully",
			"migrations_path", migrationsPath,
			"total_migrations", totalMigrations,
		)
	}

	return nil
}

// RollbackMigrations rolls back the last migration
func RollbackMigrations(migrationsPath string) error {
	if DB == nil {
		return fmt.Errorf("database connection not established")
	}

	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		slog.Info("Migrations directory does not exist, skipping rollback",
			"migrations_path", migrationsPath,
		)
		return nil
	}

	sqlDB := DB.DB
	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("failed to create MySQL driver: %w", err)
	}

	entries, err := os.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var totalRollbacks int
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		subDirPath := filepath.Join(migrationsPath, entry.Name())

		subEntries, err := os.ReadDir(subDirPath)
		if err != nil {
			slog.Warn("Failed to read subdirectory", "path", subDirPath, "error", err)
			continue
		}

		hasMigrationFiles := false
		for _, subEntry := range subEntries {
			if !subEntry.IsDir() && strings.HasSuffix(subEntry.Name(), ".up.sql") {
				hasMigrationFiles = true
				break
			}
		}

		if !hasMigrationFiles {
			continue
		}

		m, err := migrate.NewWithDatabaseInstance(
			fmt.Sprintf("file://%s", convertPathToURL(subDirPath)),
			"mysql",
			driver,
		)
		if err != nil {
			slog.Warn("Failed to create migrate instance for subdirectory",
				"path", subDirPath, "error", err)
			continue
		}

		if err := m.Steps(-1); err != nil {
			slog.Warn("Failed to rollback migrations for subdirectory",
				"path", subDirPath, "error", err)
		} else {
			totalRollbacks++
			slog.Info("Successfully rolled back migrations for subdirectory", "path", subDirPath)
		}
		m.Close()
	}

	rootEntries, err := os.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	hasRootMigrations := false
	for _, entry := range rootEntries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			hasRootMigrations = true
			break
		}
	}

	if hasRootMigrations {
		m, err := migrate.NewWithDatabaseInstance(
			fmt.Sprintf("file://%s", convertPathToURL(migrationsPath)),
			"mysql",
			driver,
		)
		if err != nil {
			return fmt.Errorf("failed to create migrate instance for root directory: %w", err)
		}
		defer m.Close()

		if err := m.Steps(-1); err != nil {
			return fmt.Errorf("failed to rollback migrations: %w", err)
		}
		totalRollbacks++
	}

	if totalRollbacks == 0 {
		slog.Info("No migration files found for rollback",
			"migrations_path", migrationsPath,
		)
	} else {
		slog.Info("Database migration rollback completed successfully",
			"migrations_path", migrationsPath,
			"total_rollbacks", totalRollbacks,
		)
	}

	return nil
}

// GetMigrationVersion returns the current migration version
func GetMigrationVersion(migrationsPath string) (uint, error) {
	if DB == nil {
		return 0, fmt.Errorf("database connection not established")
	}

	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		return 0, fmt.Errorf("migrations directory does not exist: %s", migrationsPath)
	}

	sqlDB := DB.DB
	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		return 0, fmt.Errorf("failed to create MySQL driver: %w", err)
	}

	entries, err := os.ReadDir(migrationsPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var totalVersions uint
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		subDirPath := filepath.Join(migrationsPath, entry.Name())

		subEntries, err := os.ReadDir(subDirPath)
		if err != nil {
			slog.Warn("Failed to read subdirectory", "path", subDirPath, "error", err)
			continue
		}

		hasMigrationFiles := false
		for _, subEntry := range subEntries {
			if !subEntry.IsDir() && strings.HasSuffix(subEntry.Name(), ".up.sql") {
				hasMigrationFiles = true
				break
			}
		}

		if !hasMigrationFiles {
			continue
		}

		m, err := migrate.NewWithDatabaseInstance(
			fmt.Sprintf("file://%s", convertPathToURL(subDirPath)),
			"mysql",
			driver,
		)
		if err != nil {
			slog.Warn("Failed to create migrate instance for subdirectory",
				"path", subDirPath, "error", err)
			continue
		}

		version, dirty, err := m.Version()
		if err != nil {
			slog.Warn("Failed to get migration version for subdirectory",
				"path", subDirPath, "error", err)
		} else {
			if dirty {
				return version, fmt.Errorf("database is in dirty state at version %d in subdirectory %s", version, subDirPath)
			}
			totalVersions += version
		}
		m.Close()
	}

	rootEntries, err := os.ReadDir(migrationsPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	hasRootMigrations := false
	for _, entry := range rootEntries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			hasRootMigrations = true
			break
		}
	}

	if hasRootMigrations {
		m, err := migrate.NewWithDatabaseInstance(
			fmt.Sprintf("file://%s", convertPathToURL(migrationsPath)),
			"mysql",
			driver,
		)
		if err != nil {
			return 0, fmt.Errorf("failed to create migrate instance for root directory: %w", err)
		}
		defer m.Close()

		version, dirty, err := m.Version()
		if err != nil {
			return 0, fmt.Errorf("failed to get migration version: %w", err)
		}

		if dirty {
			return version, fmt.Errorf("database is in dirty state at version %d", version)
		}

		totalVersions += version
	}

	return totalVersions, nil
}

// CreateMigration creates a new migration file
func CreateMigration(migrationsPath, name string) error {
	upFile := fmt.Sprintf("%s/%s.up.sql", migrationsPath, name)
	downFile := fmt.Sprintf("%s/%s.down.sql", migrationsPath, name)

	if err := os.WriteFile(upFile, []byte("-- Add your up migration SQL here\n"), 0644); err != nil {
		return fmt.Errorf("failed to create up migration file: %w", err)
	}

	if err := os.WriteFile(downFile, []byte("-- Add your down migration SQL here\n"), 0644); err != nil {
		return fmt.Errorf("failed to create down migration file: %w", err)
	}

	slog.Info("Created migration files",
		"name", name,
		"up_file", upFile,
		"down_file", downFile,
		"migrations_path", migrationsPath,
	)
	return nil
}
