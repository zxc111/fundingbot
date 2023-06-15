package bybit

import (
	"encoding/json"
	"github.com/zxc111/fundingbot"
	"go.uber.org/zap"
)

type (
	AccountInfo struct {
		AccountLTV    string `json:"AccountLTV"`
		AccountIMRate string `json:"accountIMRate"`
		AccountMMRate string `json:"accountMMRate"`
		TotalEquity   string `json:"totalEquity"`
	}
	Resp struct {
		RetCode int    `json:"retCode"`
		RetMsg  string `json:"retMsg"`
		Result  struct {
			List []*AccountInfo `json:"list"`
		} `json:"result"`
	}
)

func GetMMR(c fundingbot.Config, key string) {
	params := "accountType=UNIFIED"
	body, err := get(key, c.PriKey, params)
	if err != nil {
		fundingbot.Logger.Error("getMMR err", zap.Any("err", err))
		return
	}
	v := new(Resp)
	err = json.Unmarshal([]byte(body), v)
	if err != nil {
		fundingbot.Logger.Error("unmarshal err", zap.Any("err", err))
		return
	}
	fundingbot.Logger.Info("accountInfo", zap.Any("info", v.Result.List[0]))
}
