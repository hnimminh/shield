package config

import (
	"os"
	"strconv"
	"strings"
)

var (
	// webservice
	HTTPListenIP   string
	HTTPListenPort int

	// application log level
	LogLevel int8

	// Configuration Redis
	CfgRedisAddress  string
	CfgRedisPassword string
	CfgRedisDb       int
)

func init() {
	LogLevel = DFT_LOG_LEVEL
	if _level, ok := ZLOGLEVEL[strings.ToUpper(os.Getenv("LOGLEVEL"))]; ok {
		LogLevel = _level
	}

	HTTPListenIP = os.Getenv("HTTP_LISTEN_IP")
	if HTTPListenIP == "" {
		HTTPListenIP = DFT_API_HOST
	}
	HTTPListenPort, _ = strconv.Atoi(os.Getenv("HTTP_LISTEN_PORT"))
	if HTTPListenPort == 0 {
		HTTPListenPort = DFT_API_PORT
	}

	CfgRedisAddress = os.Getenv("CFG_REDIS_ADDR")
	CfgRedisPassword = os.Getenv("CFG_REDIS_PASSWORD")
	CfgRedisDb, _ = strconv.Atoi(os.Getenv("CFG_REDIS_DB"))
}
