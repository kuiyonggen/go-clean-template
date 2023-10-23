package amqprpc

import (
	"github.com/kuiyonggen/go-clean-template/internal/usecase"
	"github.com/kuiyonggen/go-clean-template/internal/usecase/repo"
	"github.com/kuiyonggen/go-clean-template/internal/usecase/webapi"
	"github.com/kuiyonggen/go-clean-template/pkg/rabbitmq/rmq_rpc/server"
        "github.com/kuiyonggen/go-clean-template/config"
    )

// NewRouter -.
func NewRouter(cfg * config.Config) map[string]server.CallHandler {
	routes := make(map[string]server.CallHandler)
	{
            // Use case
            t := usecase.New(
                repo.New(cfg.Pg),
                webapi.New(),
            )
	    newTranslationRoutes(routes, t)
	}

	return routes
}
