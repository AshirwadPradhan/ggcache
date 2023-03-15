package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/AshirwadPradhan/ggcache/cache"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
}

type Server struct {
	ServerOpts
	cache cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      c,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}
	log.Printf("server starting on port [%s]\n", s.ListenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %s\n", err)
			continue
		}
		go s.handleConn(conn)
	}

}

func (s *Server) handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
	}()

	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("conn read error: %s\n", err)
			break
		}
		go s.handleCommand(conn, buf[:n])
	}
}

func (s *Server) handleCommand(conn net.Conn, rawCmd []byte) {
	msg, err := s.parseMessage(rawCmd)
	if err != nil {
		fmt.Printf("failed to parse command %s : %s", string(rawCmd), err)
		return
	}

	switch msg.Cmd {
	case CMDSet:
		if err := s.handleSetCommand(conn, msg); err != nil {
			log.Printf("value not set [%s]\n", err)
			return
		}
	}
}

func (s *Server) handleSetCommand(conn net.Conn, msg *Message) error {
	if err := s.cache.Set(msg.Key, msg.Value, msg.TTL); err != nil {
		return err
	}

	s.sendToFollowers(context.TODO(), msg)

	return nil
}

func (s *Server) sendToFollowers(ctx context.Context, msg *Message ) error {
	return nil
}
