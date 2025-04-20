package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vterry/ddd-study/auth-server/internal/infra/api"
	"github.com/vterry/ddd-study/auth-server/internal/infra/config"
	"github.com/vterry/ddd-study/auth-server/internal/infra/db/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {

	// default context for the application
	ctx := context.Background()

	mongoUri := mongodb.MongoURIBuilder(config.Envs.MONGOAddress)
	mongoOptions := mongodb.NewMongoDBStorage(mongoUri, config.Envs.MONGOUser, config.Envs.MONGOPass)

	// another context just for initial connection

	connectCtx, connectCancel := context.WithTimeout(ctx, 5*time.Second)
	defer connectCancel()

	dbClient, err := initStorage(connectCtx, mongoOptions)
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB storage: %v", err)
	}
	defer disconnectMongoDB(dbClient, 5*time.Second)

	database := dbClient.Database("AuthServerDb")

	// another context just for server
	serverCtx, serverCancel := context.WithCancel(ctx)
	defer serverCancel()

	server := api.NewHttpServer(serverCtx, config.Envs.Port, database)
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.Run()
	}()

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		fmt.Println("Server error: %w", err)
		serverCancel() // cancel the server context
	case sig := <-shutdownChan:
		log.Printf("Received %v, shuttind down gracefully ...", sig)
		serverCancel()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		} else {
			log.Println("api server shutdown gracefully")
		}
	}
}

func initStorage(ctx context.Context, db *options.ClientOptions) (*mongo.Client, error) {

	client, err := mongo.Connect(ctx, db)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		// Clean up on connection failure
		if disconnectErr := client.Disconnect(ctx); disconnectErr != nil {
			log.Printf("Warning: failed to disconnect after connection failure: %v", disconnectErr)
		}
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	log.Println("Successfully connected to MongoDB")
	return client, nil
}

func disconnectMongoDB(client *mongo.Client, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := client.Disconnect(ctx); err != nil {
		log.Printf("MongoDB disconnect error: %v", err)
	} else {
		log.Println("MongoDB disconnected gracefully")
	}
}
