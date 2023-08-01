package main

import "github.com/iVitaliya/cache-go/framework"

// https://youtu.be/sRXIRikME14?t=1035
type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
	LeaderAddr string
}

type Server struct {
	ServerOpts,

	members map[*client.Client]struct{}

	cache framework.Cacher
}

func NewServer() {
	
}