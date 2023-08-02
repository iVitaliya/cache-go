package cache

import (
	"github.com/iVitaliya/cache-go/client"
	"github.com/iVitaliya/cache-go/framework"
)

type cacheClient struct {
	Err  error
	User *client.Client
}

type cacheServer struct {
	Server *framework.Server
	Client *cacheClient
}

func CreateDefault() *cacheServer {
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

	return &cacheServer{
		Server: server,
		Client: &cacheClient{
			Err:  err,
			User: client,
		},
	}
}

func CreateCustomServer(listenAddr string, leaderAddr string) *cacheServer {
	opts := framework.ServerOpts{
		ListenAddr: listenAddr,
		LeaderAddr: leaderAddr,
		IsLeader:   len(leaderAddr) == 0,
	}

	server := framework.NewServer(opts, framework.New())
	client, err := client.New(listenAddr, client.Options{})

	return &cacheServer{
		Server: server,
		Client: &cacheClient{
			Err:  err,
			User: client,
		},
	}
}
