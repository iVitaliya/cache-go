package cache

import (
	"github.com/iVitaliya/cache-go/framework"
)

func CreateDefaultServer() *framework.Server {
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
	return server
}

func CreateCustomServer(listenAddr string, leaderAddr string) *framework.Server {
	opts := framework.ServerOpts{
		ListenAddr: listenAddr,
		LeaderAddr: leaderAddr,
		IsLeader:   len(leaderAddr) == 0,
	}

	server := framework.NewServer(opts, framework.New())
	return server
}
