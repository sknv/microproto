package cfg

import (
	"os"

	"github.com/sknv/microproto/app/lib/xflags"
)

type Config struct {
	Addr string `long:"math-addr" env:"MATH_ADDR" default:"localhost:8081" description:"math service address"`
	// HealthAddr string `long:"math-health-addr" env:"MATH_HEALTH_ADDR" default:"localhost:8082" description:"health check server address for the math service"`
	// ConsulAddr string `long:"consul-addr" env:"CONSUL_ADDR" default:"localhost:8500" description:"consul service"`
}

func Parse() *Config {
	var cfg Config
	if _, err := xflags.ParseArgs(os.Args[1:], &cfg); err != nil {
		os.Exit(1)
	}
	return &cfg
}
