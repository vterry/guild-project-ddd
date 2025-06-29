package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/vterry/ddd-study/character/internal/infra/config"
	"github.com/vterry/ddd-study/character/internal/infra/db"
	server "github.com/vterry/ddd-study/character/internal/infra/http"
	"github.com/vterry/ddd-study/character/internal/infra/logger"
)

func main() {
	zapLogger := logger.NewZapLogger()
	zapLogger.Info("Starting Character Service", "addr", config.Envs.Addr)

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
		zapLogger.Error("failed to connect to MySQL", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := dbConn.Close(); err != nil {
			zapLogger.Error("failed to close DB connection", "error", err)
		} else {
			zapLogger.Info("DB connection closed gracefully")
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
		zapLogger.Error("HTTP server error", "error", err)
	case sig := <-shutdownChan:
		zapLogger.Info("Received shutdown signal", "signal", sig)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := httpServer.Stop(shutdownCtx); err != nil {
		zapLogger.Error("HTTP server shutdown error", "error", err)
	} else {
		zapLogger.Info("HTTP server shutdown gracefully")
	}
}
