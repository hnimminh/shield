// package main
package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hnimminh/shield/blueprint"
	"github.com/hnimminh/shield/config"
	"github.com/hnimminh/shield/versions"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

var (
	host    string
	port    int
	debug   bool
	crdburl string
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
        %sVersion v%s
    -----------------------------------------------------------------
    ` + "\n\n"
	fmt.Printf(banner, versions.NAME, versions.VERSION)

	// -------------------------------------------------------------
	flag.StringVar(&host, "host", "", "HTTP API binding IP address")
	flag.StringVar(&host, "H", "", "HTTP API binding IP address")
	flag.IntVar(&port, "port", 0, "HTTP API binding port")
	flag.IntVar(&port, "P", 0, "HTTP API binding port")
	flag.StringVar(&crdburl, "cfg-rdb-url", "", "configuation-redis database address, eg: tcp://username:password@10.10.10.10:6379/0")
	flag.StringVar(&crdburl, "c", "", "configuation-redis database address, eg: tcp://username:password@10.10.10.10:6379/0")
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
		config.HTTPListenIP = host
	}
	if port != 0 {
		config.HTTPListenPort = port
	}

	// configuration redis setting
	config.RedisCfgSettings = blueprint.RedisStruct{
		Addr:     config.CfgRedisAddress,
		Password: config.CfgRedisPassword,
		DB:       config.CfgRedisDb,
	}
	if crdburl != "" {
		u, err := url.Parse(crdburl)
		if err != nil {
			zlog.Error().Err(err).Str("function", "HmAgent:Main:Validatevar").Msg("Fail to parse cfg-rdb-url")
		} else if len(u.Path) < 2 {
			zlog.Error().Err(blueprint.ErrInvalidRedisUrl).Str("function", "HmAgent:Main:Validatevar").Msg("path is not a redis array number")
		} else {
			_cfgRedisPassword, _ := u.User.Password()
			_cfgRedisDb, _ := strconv.Atoi(u.Path[1:])
			_redisCfgSettings := blueprint.RedisStruct{
				Network:  u.Scheme,
				Addr:     u.Host,
				Username: u.User.Username(),
				Password: _cfgRedisPassword,
				DB:       _cfgRedisDb,
			}
			_crdburl := _redisCfgSettings.String()
			if crdburl == _crdburl {
				config.RedisCfgSettings = _redisCfgSettings
			} else {
				zlog.Error().Err(blueprint.ErrInvalidRedisUrl).Str("function", "HmAgent:Main:Validatevar").Msgf("url mismatch between rebuilt(%v)", _crdburl)
			}
		}
	}
}

func main() {
	basesvc.Starter(config.RedisCfgSettings.IsNone())
	if config.RedisCfgSettings.IsNone() {
		zlog.Warn().Str("function", "HmAgent:Main").Msgf("No configuration database was declared!")
		go basesvc.RdbServer.Start()
	}
	web.Server()
}
