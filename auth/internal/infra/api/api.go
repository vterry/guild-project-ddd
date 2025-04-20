package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/vterry/ddd-study/auth-server/internal/app/auth"
	"github.com/vterry/ddd-study/auth-server/internal/domain/session"
	"github.com/vterry/ddd-study/auth-server/internal/infra/db/mongodb"
	middleware "github.com/vterry/ddd-study/auth-server/internal/infra/middlware"
	"go.mongodb.org/mongo-driver/mongo"
)

type ApiServer struct {
	addr   string
	db     *mongo.Database
	server *http.Server
	ctx    context.Context
}

func NewHttpServer(ctx context.Context, addr string, db *mongo.Database) *ApiServer {
	return &ApiServer{
		addr: addr,
		db:   db,
		ctx:  ctx,
	}
}

func (a *ApiServer) Run() error {
	v1 := http.NewServeMux()
	v1.Handle("/authserver/v1/", middleware.Chain(http.StripPrefix("/authserver/v1", v1), middleware.Logger()))

	loginRepository, err := mongodb.NewLoginRepository(a.ctx, a.db)
	if err != nil {
		return fmt.Errorf("erro while initiatilizing Login Repo")
	}

	sessionRepository, err := mongodb.NewSessionRepository(a.ctx, a.db)
	if err != nil {
		return fmt.Errorf("erro while initiatilizing Login Repo")
	}

	sessionService := session.NewSessionService(sessionRepository, loginRepository)

	authService := auth.NewAuthService(sessionService, loginRepository)
	handler := auth.NewHandler(*authService)
	handler.RegisterRoutes(v1)

	a.server = &http.Server{
		Addr:    a.addr,
		Handler: v1,
	}

	log.Println("Listening on", a.addr)
	return a.server.ListenAndServe()
}

func (a *ApiServer) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
