// Package v1 implements routing paths. Each services in own file.
package v1

import (
        "fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Swagger docs.
	"github.com/kuiyonggen/go-clean-template/docs"
        "github.com/kuiyonggen/go-clean-template/config"
    )

// NewRouter _.
func NewRouter(handler *gin.Engine, cfg *config.Config) {
        basePath := "/v1"
        // Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
        // programmatically set swagger info
        if cfg.Swagger {
            docs.SwaggerInfo.Title = fmt.Sprintf("%s Service API", cfg.Name)
            docs.SwaggerInfo.Description = fmt.Sprintf("%s Service.", cfg.Name)
            docs.SwaggerInfo.Version = cfg.Version
            docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.Address, cfg.Port)
            docs.SwaggerInfo.BasePath = basePath
	    swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	    handler.GET("/swagger/*any", swaggerHandler)
        }
	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	h := handler.Group(basePath)
	{
	    newTranslationRoutes(h, cfg)
            newHelloRoutes(h, cfg)
	}
}
