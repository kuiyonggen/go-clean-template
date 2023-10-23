package config

import (
	"fmt"
        "flag"
        "strings"
	"github.com/ilyakaznacheev/cleanenv"
        "github.com/kuiyonggen/go-clean-template/pkg/consul"    
        "github.com/kuiyonggen/go-clean-template/pkg/postgres"
        "github.com/kuiyonggen/go-clean-template/pkg/logger"
        "github.com/kuiyonggen/go-clean-template/utils"
    )

type (
	// Config -.
	Config struct {
		App  `yaml:"app"`
                Log    `yaml:"log"`
		HTTP `yaml:"http"`
                Consul `yaml:"consul"`
		PG   `yaml:"postgres"`
		RMQ  `yaml:"rabbitmq"`
                Extra map[string]interface{} `yaml:"extra"`

                CClient *consul.Consul
                Pg *postgres.Postgres
                Logger *logger.Logger
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Address string `env-required:"true" yaml:"address" env:"HTTP_ADDRESS"`
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

        // Log -.
        Log struct {
                Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
        }

        // Consul -.
        Consul struct {
            CheckApi string `yaml:"check_api"`
            Interval string `yaml:"interval"`
            Timeout string  `yaml:"timeout"`
            Tags []string `yaml:"tags"`
        }

	// PG -.
	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     string `env-required:"true"                 env:"PG_URL"`
	}

	// RMQ -.
	RMQ struct {
		ServerExchange string `env-required:"false" yaml:"rpc_server_exchange" env:"RMQ_RPC_SERVER"`
		ClientExchange string `env-required:"false" yaml:"rpc_client_exchange" env:"RMQ_RPC_CLIENT"`
		URL            string `env-required:"false"                            env:"RMQ_URL"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	// config args, priority: config > consul
	var (
		configFile   = flag.String("config", "", "config file, prior to use.")
		consulAddr   = flag.String("consul", "localhost:8500", "consul server address.")
		consulFolder = flag.String("folder", "", "consul kv folder.")
		serviceName  = flag.String("name", "microapp", "both microservice name and kv name.")
		listenAddr   = flag.String("listen", ":8080", "listen address.")
	)
	flag.Parse()

	cfg := &Config{}

	if len(*configFile) == 0 {
		cfg.Name = *serviceName
		cfg.Version = "v0.0.1"
                la := strings.Split(*listenAddr, ":")
                if len(la[0]) == 0 {
                    cfg.HTTP.Address = utils.GetHostIP()
                } else {
                    cfg.HTTP.Address = la[0]
                }
                cfg.HTTP.Port = la[1]
                c, err := consul.New(*consulAddr, *consulFolder, *serviceName)
                cfg.CClient = c
                if err != nil {
                    fmt.Printf("failed to create consul endpoint, error: %v.\n", err)
                }
                err = cfg.CClient.Kv(cfg)
	} else {
		err := cleanenv.ReadConfig(*configFile, cfg)
		if err != nil {
			return nil, fmt.Errorf("config error: %w", err)
		}

		err = cleanenv.ReadEnv(cfg)
		if err != nil {
			return nil, err
		}
	}

	return cfg, nil
}
