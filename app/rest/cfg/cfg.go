package cfg

import (
	"os"

	"github.com/sknv/microproto/app/lib/xflags"
)

type Config struct {
	Addr          string `long:"rest-addr" env:"REST_ADDR" default:":8080" description:"rest api address"`
	ConsulAddr    string `long:"consul-addr" env:"CONSUL_ADDR" default:"127.0.0.1:8500" description:"consul address"`
	MathProxyAddr string `long:"math-proxy-addr" env:"MATH_PROXY_ADDR" default:"http://127.0.0.1:8001" description:"proxy address for math services"`
}

func Parse() *Config {
	cfg := new(Config)
	if _, err := xflags.ParseArgs(os.Args[1:], cfg); err != nil {
		os.Exit(1)
	}
	return cfg
}
