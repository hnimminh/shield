// package main
package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hnimminh/shield/internal/blueprint"
	"github.com/hnimminh/shield/internal/config"
	"github.com/hnimminh/shield/pubsub"
	"github.com/hnimminh/shield/web"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

var (
	host     string
	port     int
	redisurl string
	channel  string
	debug    bool
)

func init() {
	// -------------------------------------------------------------
	// startup banner with setting displayed
	// -------------------------------------------------------------
	banner := `
    +---------------------------------------------------------------+
	|	 __  __  __ __  ____ __    ____								|
	|	(( \ ||  || || ||    ||    || \\							|
	|	 \\  ||==|| || ||==  ||    ||  ))							|
	|	\_)) ||  || || ||___ ||__| ||_//							|
	|																|
    +---------------------------------------------------------------+
        Simple Daemon receiving/executing firewall config command
        v%s
    -----------------------------------------------------------------
    ` + "\n\n"
	fmt.Printf(banner, config.Version)

	// -------------------------------------------------------------
	flag.StringVar(&host, "host", "", "HTTP API binding IP address")
	flag.StringVar(&host, "H", "", "HTTP API binding IP address")
	flag.IntVar(&port, "port", 0, "HTTP API binding port")
	flag.IntVar(&port, "P", 0, "HTTP API binding port")
	flag.StringVar(&redisurl, "redisurl", "", "redis url, eg: tcp://username:password@10.10.10.10:6379/0")
	flag.StringVar(&redisurl, "r", "", "redis url, eg: tcp://username:password@10.10.10.10:6379/0")
	flag.StringVar(&channel, "channel", "", "redis channel for pubsub")
	flag.StringVar(&channel, "c", "", "redis channel for pubsub")
	flag.BoolVar(&debug, "debug", false, "sets log level to debug")
	flag.BoolVar(&debug, "d", false, "sets log level to debug")
	flag.Parse()

	// -------------------------------------------------------------
	// log setting
	// -------------------------------------------------------------
	if debug && config.LogLevel > 0 {
		config.LogLevel = 0
	}
	zerolog.SetGlobalLevel(zerolog.Level(config.LogLevel))
	zlog.Logger = zlog.Output(
		zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			FormatLevel: func(i interface{}) string {
				return strings.ToUpper(fmt.Sprintf("[%4s]", i))
			},
			NoColor: false,
		},
	)

	//-------------------------------------------------------------
	//	variable validation
	//	dirrect value (via cli) will overide envar value
	//-------------------------------------------------------------
	// http listen service
	if host != "" {
		zlog.Info().Str("function", "Shield:Main:Validatevar").Msgf("Overwrite host to %s", host)
		config.HTTPListenIP = host
	}
	if port != 0 {
		zlog.Info().Str("function", "Shield:Main:Validatevar").Msgf("Overwrite port to %d", port)
		config.HTTPListenPort = port
	}

	// configuration redis setting
	config.RedisCfgSettings = blueprint.RedisStruct{
		Addr:     config.RedisAddress,
		Password: config.RedisPassword,
		DB:       config.RedisDb,
	}
	if redisurl != "" {
		u, err := url.Parse(redisurl)
		if err != nil {
			zlog.Error().Err(err).Str("function", "Shield:Main:Validatevar").Msg("Fail to parse redis url")
		} else if len(u.Path) < 2 {
			zlog.Error().Err(blueprint.ErrInvalidRedisUrl).Str("function", "Shield:Main:Validatevar").Msg("Path is not a redis array number")
		} else {
			_redisPassword, _ := u.User.Password()
			_redisDb, _ := strconv.Atoi(u.Path[1:])
			_redisCfgSettings := blueprint.RedisStruct{
				Network:  u.Scheme,
				Addr:     u.Host,
				Username: u.User.Username(),
				Password: _redisPassword,
				DB:       _redisDb,
			}
			_redisurl := _redisCfgSettings.String()
			if redisurl == _redisurl {
				config.RedisCfgSettings = _redisCfgSettings
			} else {
				zlog.Error().Err(blueprint.ErrInvalidRedisUrl).Str("function", "Shield:Main:Validatevar").Msgf("URL mismatch with reversed(%v)", _redisurl)
			}
		}
	}
	// redis pubsub
	if channel != "" {
		zlog.Info().Str("function", "Shield:Main:Validatevar").Msgf("Overwrite pubsub channel name to %s", channel)
		config.RedisPubsubChannel = channel
	}
}

func main() {
	if (config.HTTPListenIP == "" && config.HTTPListenPort == 0) || config.RedisCfgSettings.IsNone() {
		zlog.Error().Str("function", "Shield:Main").Msgf("Use at least one of these method `API` or `PUBSUB`")
		os.Exit(1)
	}

	wg := sync.WaitGroup{}
	defer wg.Wait()

	if config.HTTPListenIP != "" && config.HTTPListenPort != 0 {
		zlog.Info().Str("function", "Shield:Main").Msgf("Listen command via HTTP API server")
		wg.Add(1)
		go web.Server(&wg)
	}

	if !config.RedisCfgSettings.IsNone() {
		zlog.Info().Str("function", "Shield:Main").Msgf("Listen command via Redis Pub/Sub")
		wg.Add(1)
		go pubsub.Eventd(&wg)
	}
}
