package http

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/vterry/ddd-study/character/internal/adapters/input/rest"
	"github.com/vterry/ddd-study/character/internal/adapters/output/gateway"
	"github.com/vterry/ddd-study/character/internal/adapters/output/repository/mysql"
	"github.com/vterry/ddd-study/character/internal/core/service"
	"github.com/vterry/ddd-study/character/internal/infra/config"
	"github.com/vterry/ddd-study/character/internal/infra/keycloak"
)

type HttpServer struct {
	addr   string
	db     *sql.DB
	server *http.Server
	ctx    context.Context
}

func NewHttpServer(ctx context.Context, addr string, db *sql.DB) *HttpServer {
	return &HttpServer{
		addr: addr,
		db:   db,
		ctx:  ctx,
	}
}

func (h *HttpServer) Run() error {

	characterRepo := mysql.NewCharacterRepository(h.db)

	vaultGateway := gateway.NewMockVaultGateway()

	keycloakClient, err := keycloak.NewKeycloakClient(h.ctx, &config.Envs.Auth)
	if err != nil {
		return err
	}
	loginGateway := gateway.NewLoginGateway(keycloakClient)
	loginGateway.Client = &http.Client{}

	characterCoreService := service.NewCharacterService(characterRepo, vaultGateway)

	characterService := rest.NewCharacterService(characterCoreService, loginGateway)

	handler := rest.NewHandler(*characterService)

	v1 := http.NewServeMux()
	v1.Handle("/character/v1/", http.StripPrefix("/character/v1", v1))
	handler.RegisterRoutes(v1)

	h.server = &http.Server{
		Addr:    h.addr,
		Handler: v1,
	}

	return h.server.ListenAndServe()
}

func (h *HttpServer) Stop(ctx context.Context) error {
	return h.server.Shutdown(ctx)
}
