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
	RedisAddress       string
	RedisPassword      string
	RedisDb            int
	RedisPubsubChannel string
)

func init() {
	LogLevel = DFT_LOG_LEVEL
	if _level, ok := ZLOGLEVEL[strings.ToUpper(os.Getenv("LOGLEVEL"))]; ok {
		LogLevel = _level
	}

	// http api
	HTTPListenIP = os.Getenv("HTTP_LISTEN_IP")
	HTTPListenPort, _ = strconv.Atoi(os.Getenv("HTTP_LISTEN_PORT"))

	// redis address
	RedisAddress = os.Getenv("REDIS_ADDR")
	RedisPassword = os.Getenv("REDIS_PASSWORD")
	RedisDb, _ = strconv.Atoi(os.Getenv("REDIS_DB"))
	// redis pubsub channel
	RedisPubsubChannel = os.Getenv("REDIS_PUBSUB_CHANNEL")
	if RedisPubsubChannel == "" {
		RedisPubsubChannel = "SECURITY_CHANNEL"
	}

}
