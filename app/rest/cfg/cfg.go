package cfg

import (
	"os"

	flags "github.com/jessevdk/go-flags"
)

type Config struct {
	Addr string `long:"rest-addr" env:"REST_ADDR" default:"localhost:8080" description:"rest api address"`
}

func Parse() *Config {
	config := &Config{}
	flagParser := flags.NewParser(config, flags.Default)
	if _, err := flagParser.ParseArgs(os.Args[1:]); err != nil {
		os.Exit(1)
	}
	return config
}
