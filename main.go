package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/iVitaliya/cache-go/client"
	"github.com/iVitaliya/cache-go/framework"
	"github.com/iVitaliya/logger-go"
)

func main() {
	var (
		listenAddr = flag.String("listenaddr", ":3000", "listen address of the server")
		leaderAddr = flag.String("leaderaddr", "", "listen address of the leader")
	)
	flag.Parse()

	opts := ServerOpts{
		ListenAddr: *listenAddr,
		IsLeader:   len(*leaderAddr) == 0,
		LeaderAddr: *leaderAddr,
	}

	go func() {
		time.Sleep(time.Second * 10)
		if opts.IsLeader {
			SendStuff()
		}
	}()

	server := NewServer(opts, framework.New())
	server.Start()
}

func SendStuff() {
	for i := 0; i < 100; i++ {
		go func(i int) {
			client, err := client.New(":3000", client.Options{})
			if err != nil {
				logger.Error(err)
			}

			var (
				key   = []byte(fmt.Sprintf("key_%d", i))
				value = []byte(fmt.Sprintf("value_%d", i))
			)

			err = client.Set(context.Background(), key, value, 0)
			if err != nil {
				logger.Error(err)
			}

			fetchedValue, err := client.Get(context.Background(), key)
			if err != nil {
				logger.Error(err)
			}
			logger.Debug(fetchedValue)

			client.Close()
		}(i)
	}
}
