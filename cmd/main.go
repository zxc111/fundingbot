package main

import (
	"github.com/zxc111/fundingbot"
	"github.com/zxc111/fundingbot/internal/bybit"
)

func main() {
	fundingbot.InitConfig()
	bybit.InitHttpClient(fundingbot.C)
	for _, k := range fundingbot.C.Keys {
		bybit.GetMMR(fundingbot.C, k)
	}
}
