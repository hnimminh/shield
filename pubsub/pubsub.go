package pubsub

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hnimminh/shield/internal/config"
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

		zlog.Info().Str("function", "Shield:Pubsub:Eventd").Msgf("popout channel=%s, payload %v; %s %v", msg.Channel, msg.Payload, msg.Pattern, msg.PayloadSlice)
		go handler()

	}
	// marking groutine finish
	wg.Done()
}

func handler() {
	zlog.Info().Str("function", "Shield:Pubsub:Handler").Msg("just a message")
}

/*
func libreHandler() {
	zlog.Info().Str("function", "Shield:Pubsub:LibreHandler").Msg("just a message")
}
*/
