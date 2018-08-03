package cfg

import (
	"os"

	"github.com/sknv/microproto/app/lib/xflags"
)

type Config struct {
	Addr     string `long:"rest-addr" env:"REST_ADDR" default:"localhost:8080" description:"rest api address"`
	MathAddr string `long:"math-addr" env:"MATH_ADDR" default:"localhost:8081" description:"math service address"`
	// ConsulAddr string `long:"consul-addr" env:"CONSUL_ADDR" default:"localhost:8500" description:"consul service"`
}

func Parse() *Config {
	cfg := new(Config)
	if _, err := xflags.ParseArgs(os.Args[1:], cfg); err != nil {
		os.Exit(1)
	}
	return cfg
}
