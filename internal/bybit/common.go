package bybit

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/zxc111/fundingbot"
)

const (
	bybitURL  = "https://api.bybit.com"
	walletURI = "/v5/account/wallet-balance"

	window = 5000
)

var (
	Client *http.Client
)

func InitHttpClient(config fundingbot.Config) {
	tp := &http.Transport{
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	if len(config.ProxyPath) != 0 {
		u, _ := url.Parse(config.ProxyPath)
		tp.Proxy = http.ProxyURL(u)
	}
	Client = &http.Client{Transport: tp}
}

func get(key, priKey, params string) (string, error) {
	ts := int(time.Now().UnixMilli())
	sign := getSign(ts, key, window, params, priKey)
	request, err := http.NewRequest("GET", bybitURL+walletURI+"?"+params, nil)
	if err != nil {
		panic(err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-BAPI-API-KEY", key)
	request.Header.Set("X-BAPI-SIGN", sign)
	request.Header.Set("X-BAPI-TIMESTAMP", strconv.Itoa(ts))
	request.Header.Set("X-BAPI-SIGN-TYPE", "2")
	request.Header.Set("X-BAPI-RECV-WINDOW", strconv.Itoa(window))
	resp, err := Client.Do(request)
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return string(body), nil
}

func getSign(timestamp int, apiKey string, recvWindows int, queryString, priKey string) string {
	buf := bytes.Buffer{}
	buf.WriteString(strconv.Itoa(timestamp))
	buf.WriteString(apiKey)
	buf.WriteString(strconv.Itoa(recvWindows))
	buf.WriteString(queryString)
	hashed := sha256.Sum256(buf.Bytes())
	block, _ := pem.Decode([]byte(priKey))

	// 解析 RSA 私钥
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(signature)
}
