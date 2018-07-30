package cfg

import (
	"os"

	flags "github.com/jessevdk/go-flags"
)

type Config struct {
	Addr       string `long:"math-addr" env:"MATH_ADDR" default:"localhost:8081" description:"math service address"`
	HealthAddr string `long:"math-health-addr" env:"MATH_HEALTH_ADDR" default:"localhost:8082" description:"health check server address for the math service"`
	ConsulAddr string `long:"consul-addr" env:"CONSUL_ADDR" default:"localhost:8500" description:"consul service"`
}

func Parse() *Config {
	var cfg Config
	flagParser := flags.NewParser(&cfg, flags.Default)
	if _, err := flagParser.ParseArgs(os.Args[1:]); err != nil {
		os.Exit(1)
	}
	return &cfg
}
