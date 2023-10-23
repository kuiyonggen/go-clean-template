// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/kuiyonggen/go-clean-template/config"
	amqprpc "github.com/kuiyonggen/go-clean-template/internal/controller/amqp_rpc"
	v1 "github.com/kuiyonggen/go-clean-template/internal/controller/http/v1"
	"github.com/kuiyonggen/go-clean-template/pkg/httpserver"
	"github.com/kuiyonggen/go-clean-template/pkg/logger"
	"github.com/kuiyonggen/go-clean-template/pkg/postgres"
	"github.com/kuiyonggen/go-clean-template/pkg/rabbitmq/rmq_rpc/server"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
        cfg.Logger = l

        // Register Service
        serviceID, err := cfg.CClient.Register(cfg.HTTP.Address, cfg.HTTP.Port, cfg.CheckApi, 
                cfg.Interval, cfg.Timeout, cfg.Tags)
        if err != nil {
            l.Fatal(fmt.Errorf("app - Run - consul register: %w.", err))
        }
        defer cfg.CClient.Deregister(serviceID)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

        cfg.Pg = pg

	// RabbitMQ RPC Server
        rmqRouter := amqprpc.NewRouter(cfg)

        rmqServer, err := server.New(cfg.RMQ.URL, cfg.RMQ.ServerExchange, rmqRouter, l)
	if err != nil {
	    l.Fatal(fmt.Errorf("app - Run - rmqServer - server.New: %w", err))
	}

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, cfg)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case err = <-rmqServer.Notify():
		l.Error(fmt.Errorf("app - Run - rmqServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	err = rmqServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - rmqServer.Shutdown: %w", err))
	}
}
