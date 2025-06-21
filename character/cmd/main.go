package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/vterry/ddd-study/character/internal/infra/config"
	"github.com/vterry/ddd-study/character/internal/infra/db"
	server "github.com/vterry/ddd-study/character/internal/infra/http"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mysqlCfg := mysql.Config{
		User:                 config.Envs.Db.User,
		Passwd:               config.Envs.Db.Password,
		Addr:                 config.Envs.Db.Address,
		DBName:               config.Envs.Db.Name,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	dbConn, err := db.NewMySQLStorage(mysqlCfg)
	if err != nil {
		log.Fatalf("failed to connect to MySQL: %v", err)
	}
	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Printf("failed to close DB connection: %v", err)
		} else {
			log.Println("DB connection closed gracefully")
		}
	}()

	httpServer := server.NewHttpServer(ctx, config.Envs.Addr, dbConn)

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- httpServer.Run()
	}()

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		log.Printf("HTTP server error: %v", err)
	case sig := <-shutdownChan:
		log.Printf("Received signal %v, shutting down gracefully...", sig)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := httpServer.Stop(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	} else {
		log.Println("HTTP server shutdown gracefully")
	}
}
