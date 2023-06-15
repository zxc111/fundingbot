package fundingbot

import (
	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
)

type Config struct {
	PriKey    string
	Keys      []string
	ProxyPath string
}

var (
	C      Config
	Logger *zap.Logger
)

func InitConfig() {
	_, err := toml.DecodeFile("config.toml", &C)
	if err != nil {
		panic(err)
	}
	Logger, _ = zap.NewDevelopment()
}
