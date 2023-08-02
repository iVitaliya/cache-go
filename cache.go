package cache

import (
	"github.com/iVitaliya/cache-go/client"
	"github.com/iVitaliya/cache-go/framework"
)

type CacheClient struct {
	Err  error
	User *client.Client
}

type CacheServer struct {
	Server *framework.Server
	Client *CacheClient
}

func CreateDefault() *CacheServer {
	var (
		listenAddr = ":3000"
		leaderAddr = ""
	)

	opts := framework.ServerOpts{
		ListenAddr: listenAddr,
		IsLeader:   true,
		LeaderAddr: leaderAddr,
	}

	server := framework.NewServer(opts, framework.New())
	client, err := client.New(listenAddr, client.Options{})

	return &CacheServer{
		Server: server,
		Client: &CacheClient{
			Err:  err,
			User: client,
		},
	}
}

func CreateCustomServer(listenAddr string, leaderAddr string) *CacheServer {
	opts := framework.ServerOpts{
		ListenAddr: listenAddr,
		LeaderAddr: leaderAddr,
		IsLeader:   len(leaderAddr) == 0,
	}

	server := framework.NewServer(opts, framework.New())
	client, err := client.New(listenAddr, client.Options{})

	return &CacheServer{
		Server: server,
		Client: &CacheClient{
			Err:  err,
			User: client,
		},
	}
}
