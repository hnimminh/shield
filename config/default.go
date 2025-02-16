package config

import (
	"github.com/hnimminh/shield/blueprint"
)

const (
	DFT_API_HOST  = "0.0.0.0"
	DFT_API_PORT  = 8088
	DFT_LOG_LEVEL = 1
)

var (
	ZLOGLEVEL = map[string]int8{
		"TRACE":    -1,
		"DEBUG":    0,
		"INFO":     1,
		"WARNING":  2,
		"ERROR":    3,
		"CRITICAL": 4,
	}

	RedisCfgSettings blueprint.RedisStruct
)
