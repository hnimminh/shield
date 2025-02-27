package pubsub

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hnimminh/shield/internal/config"
	"github.com/hnimminh/shield/internal/nftcli"
	zlog "github.com/rs/zerolog/log"
)

var (
	ctx context.Context
	rdb *redis.Client
)

func Eventd(wg *sync.WaitGroup) {
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Network:  config.RedisCfgSettings.Network,
		Addr:     config.RedisCfgSettings.Addr,
		Password: config.RedisCfgSettings.Password,
		DB:       config.RedisCfgSettings.DB,
	})

	subscriber := rdb.Subscribe(ctx, config.RedisPubsubChannel)
	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil && err != redis.Nil {
			zlog.Error().Err(err).Str("function", "Shield:Pubsub:Eventd").Msg("Fail to subcribe channel, take a rest a few second before retry")
			time.Sleep(5 * time.Second)
			continue
		}

		zlog.Debug().Str("function", "Shield:Pubsub:Eventd").Msgf("%s popout channel=%s, payload %v; %s %v", config.NodeID, msg.Channel, msg.Payload, msg.Pattern, msg.PayloadSlice)
		go handler(msg.Payload)

	}

	// marking groutine finish
	wg.Done()
}

type FwEvent struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func handler(payload string) {
	var fe FwEvent
	if err := json.Unmarshal([]byte(payload), &fe); err != nil {
		zlog.Error().Err(err).Str("function", "Shield:Pubsub:Handler").Msgf("Unable to unmarshal data %s", payload)
		return
	}

	if fe.Type == "stream" {
		if err := nftcli.StreamNftOverShell(fe.Data); err != nil {
			zlog.Error().Err(err).Str("function", "Shield:Pubsub:Handler").Msgf("nftables can not execute stream")
		} else {
			zlog.Info().Str("function", "Shield:Pubsub:Handler").Msgf("nftables stream updated")
		}
	}

	if fe.Type == "cmd" {
		if _, err := nftcli.NftOverShell(fe.Data); err != nil {
			zlog.Error().Err(err).Str("function", "Shield:Pubsub:Handler").Msgf("nftables can not execute command %s", fe.Data)
		} else {
			zlog.Info().Str("function", "Shield:Pubsub:Handler").Msgf("nftables cmd `%s` executed", fe.Data)
		}
	}
}
