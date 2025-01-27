package kyg


import (
	"log"

	"github.com/kuiyonggen/go-clean-template/config"
	"github.com/kuiyonggen/go-clean-template/internal/app"
)

func Run () {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}

