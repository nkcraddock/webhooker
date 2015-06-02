package main

import (
	"log"

	"github.com/michaelklishin/rabbit-hole"
	"github.com/nkcraddock/webhooker/db"
	"github.com/nkcraddock/webhooker/q"
	"github.com/nkcraddock/webhooker/webhooks"
	"gopkg.in/redis.v3"
)

func main() {
	opts := getOpts()

	// Get a rabbit connection
	rab := connectRabbit(opts)

	// Get a redis connection
	red := connectRedis(opts)

	// get all the hooks from rabbit
	rabbitHooks, _ := rab.GetHooks("")

	// Remove unknown hooks from rabbit
	for _, h := range rabbitHooks {
		if _, err := red.GetHook(h.Id); err != nil {
			log.Println("DELETE", h.Id)
		}
	}

	// Get all the hooks from redis
	redisHooks, _ := red.GetHooks("")

	// Add missing hooks to rabbit
	for _, h := range redisHooks {
		if _, err := rab.GetHook(h.Id); err != nil {
			rab.SaveHook(h)
			log.Println("ADD HOOK", h.Id)
		}

		// get all the filters from rabbit
		rabbitFilters, _ := rab.GetFilters(h.Id)

		// Remove unknown filters from rabbit
		for _, f := range rabbitFilters {
			if _, err := red.GetFilter(h.Id, f.Id); err != nil {
				log.Println("DELETE FILTER", f.Id)
			}
		}

		// Get all the filters from redis
		redisFilters, _ := red.GetFilters(h.Id)

		// Add missing filters to rabbit
		for _, f := range redisFilters {
			if _, err := rab.GetFilter(h.Id, f.Id); err != nil {
				rab.SaveFilter(f)
				red.SaveFilter(f) // update redis so it gets the RabbitMQ PropsKey
				log.Println("ADD FILTER", f.Id)
			}
		}

	}

}

func connectRabbit(opts options) *q.RabbitStore {
	rmq, err := rabbithole.NewClient(opts.RabbitUri, opts.RabbitUser, opts.RabbitPass)
	if err != nil {
		panic(err)
	}

	return q.NewRabbitStore(rmq, opts.RabbitVhost)
}

func connectRedis(opts options) webhooks.Store {
	return db.RedisHookerStore(func() *redis.Client {
		return redis.NewClient(&redis.Options{
			Addr:     opts.RedisUri,
			PoolSize: 5,
			DB:       int64(opts.RedisDb),
		})
	})
}
