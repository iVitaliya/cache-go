package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/iVitaliya/cache-go/client"
	"github.com/iVitaliya/cache-go/framework"
	"github.com/iVitaliya/cache-go/protocol"
	"github.com/iVitaliya/logger-go"
)

// https://youtu.be/sRXIRikME14?t=1035
type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
	LeaderAddr string
}

type Server struct {
	ServerOpts

	members map[*client.Client]struct{}

	cache framework.Cacher
}

func NewServer(opts ServerOpts, c framework.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      c,
		members:    make(map[*client.Client]struct{}),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}

	if !s.IsLeader && len(s.LeaderAddr) != 0 {
		go func() {
			if err := s.dialLeader(); err != nil {
				logger.Error(err)
			}
		}()
	}

	logger.Debug(fmt.Sprintf("server starting on port [%s]", s.ListenAddr))

	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.Error(fmt.Sprintf("accept error: %s", err))
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *Server) dialLeader() error {
	conn, err := net.Dial("tcp", s.LeaderAddr)
	if err != nil {
		return fmt.Errorf("failed to dial leader [%s]", s.LeaderAddr)
	}

	logger.Debug(fmt.Sprintf("connected to leader: %s", s.LeaderAddr))

	binary.Write(conn, binary.LittleEndian, protocol.CmdJoin)

	s.handleConn(conn)

	return nil
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	for {
		cmd, err := protocol.ParseCommand(conn)

		if err != nil {
			if err == io.EOF {
				break
			}

			logger.Error(fmt.Sprintf("parse command error: %s", err))
			break
		}

		go s.handleCommand(conn, cmd)
	}
}

func (s *Server) handleCommand(conn net.Conn, cmd any) {
	switch v := cmd.(type) {
	case *protocol.CommandSet:
		s.handleSetCommand(conn, v)
	case *protocol.CommandGet:
		s.handleGetCommand(conn, v)
	case *protocol.CommandJoin:
		s.handleJoinCommand(conn, v)
	}
}

func (s *Server) handleJoinCommand(conn net.Conn, cmd *protocol.CommandJoin) error {
	logger.Debug(fmt.Sprintf("member just joined the cluster: ", conn.RemoteAddr()))

	s.members[client.NewFromConn(conn)] = struct{}{}

	return nil
}

func (s *Server) handleGetCommand(conn net.Conn, cmd *protocol.CommandGet) error {
	resp := protocol.ResponseGet{}
	value, err := s.cache.Get(cmd.Key)
	if err != nil {
		resp.Status = protocol.StatusError
		_, err := conn.Write(resp.Bytes())
		return err
	}

	resp.Status = protocol.StatusOK
	resp.Value = value
	_, err = conn.Write(resp.Bytes())

	return err
}

func (s *Server) handleSetCommand(conn net.Conn, cmd *protocol.CommandSet) error {
	logger.Debug(fmt.Sprintf("SET \"%s\" to \"%s\"", cmd.Key, cmd.Value))

	go func() {
		for member := range s.members {
			err := member.Set(context.TODO(), cmd.Key, cmd.Value, cmd.TTL)
			if err != nil {
				logger.Error(fmt.Sprintf("forward to member error: %s", err))
			}
		}
	}()

	resp := protocol.ResponseSet{}
	if err := s.cache.Set(cmd.Key, cmd.Value, time.Duration(cmd.TTL)); err != nil {
		resp.Status = protocol.StatusError
		_, err := conn.Write(resp.Bytes())

		return err
	}

	resp.Status = protocol.StatusOK
	_, err := conn.Write(resp.Bytes())

	return err
}
