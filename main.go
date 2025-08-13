package main

import (
	"fmt"
	"log"
	"os"

	configure "circle-center/globals/configure"
	dbpkg "circle-center/globals/db"
	"circle-center/globals/mail"
	"circle-center/panel/account"
	accountsvc "circle-center/panel/account/svc"
	mgr "circle-center/panel/manager"
	editor "circle-center/processor"
	"circle-center/reader"
)

func main() {
	cfg, err := configure.LoadConfig("")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if err := dbpkg.ConnectMySQL(&dbpkg.MySQLConfig{
		Host:           cfg.MySQL.Host,
		Port:           cfg.MySQL.Port,
		Username:       cfg.MySQL.Username,
		Password:       cfg.MySQL.Password,
		Database:       cfg.MySQL.Database,
		Charset:        cfg.MySQL.Charset,
		ParseTime:      cfg.MySQL.ParseTime,
		Loc:            cfg.MySQL.Loc,
		MaxOpenConns:   cfg.MySQL.MaxOpenConns,
		MaxIdleConns:   cfg.MySQL.MaxIdleConns,
		MaxLifetime:    cfg.MySQL.MaxLifetime,
		MultiStatement: cfg.MySQL.MultiStatement,
	}); err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer dbpkg.CloseMySQL()

	if err := dbpkg.ConnectRedis(&dbpkg.RedisConfig{
		Host:         cfg.Redis.Host,
		Port:         cfg.Redis.Port,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
		MaxRetries:   cfg.Redis.MaxRetries,
		DialTimeout:  cfg.Redis.DialTimeout,
		ReadTimeout:  cfg.Redis.ReadTimeout,
		WriteTimeout: cfg.Redis.WriteTimeout,
		IdleTimeout:  cfg.Redis.IdleTimeout,
	}); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer dbpkg.CloseRedis()

	if err := configure.InitializeJWTKeys(&cfg.JWT); err != nil {
		log.Fatalf("Failed to initialize JWT keys: %v", err)
	}

	mailService, err := mail.NewMailService(&mail.MailConfig{
		Host:     cfg.Mail.Host,
		Port:     cfg.Mail.Port,
		Username: cfg.Mail.Username,
		Password: cfg.Mail.Password,
		From:     cfg.Mail.From,
		TLSMode:  cfg.Mail.TLSMode,
	})
	if err != nil {
		log.Fatalf("Failed to initialize mail service: %v", err)
	}
	defer mailService.Close()

	if _, err := os.Stat("repository/migrations"); err == nil {
		if err := dbpkg.RunMigrations("repository/migrations"); err != nil {
			log.Printf("Warning: Failed to run migrations: %v", err)
		}
	}

	jwtClient, err := accountsvc.NewJWTClientFromGlobalKeys()
	if err != nil {
		log.Fatalf("Failed to create JWT client from global keys: %v", err)
	}
	authClient := accountsvc.NewAuthClient(jwtClient)

	r := configure.SetupRouter()

	v1 := r.Group("/v1")
	reader.RegisterRoutes(v1)
	editor.RegisterRoutes(v1)
	account.RegisterRoutes(v1, dbpkg.GetDB().DB, mailService, authClient)
	mgr.RegisterRoutes(v1, dbpkg.GetDB().DB, authClient)

	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
