package cfg

import (
	"os"

	flags "github.com/jessevdk/go-flags"
)

type Config struct {
	Addr    string `long:"rest-addr" env:"REST_ADDR" default:"localhost:8080" description:"rest api address"`
	MathURL string `long:"math-url" env:"MATH_URL" default:"http://localhost:8081" description:"math service url"`
	// ConsulAddr string `long:"consul-addr" env:"CONSUL_ADDR" default:"localhost:8500" description:"consul service"`
}

func Parse() *Config {
	var cfg Config
	flagParser := flags.NewParser(&cfg, flags.Default)
	if _, err := flagParser.ParseArgs(os.Args[1:]); err != nil {
		os.Exit(1)
	}
	return &cfg
}
