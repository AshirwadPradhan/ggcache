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
		return fmt.Errorf("main: listen error: %s", err)
	}
	log.Printf("main: server starting on port [%s]\n", s.ListenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("server: accept error: %s\n", err)
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
			log.Printf("server: conn read error: %s\n", err)
			break
		}
		go s.handleCommand(conn, buf[:n])
	}
}

func (s *Server) handleCommand(conn net.Conn, rawCmd []byte) {
	msg, err := s.parseMessage(rawCmd)
	if err != nil {
		fmt.Printf("server: failed to parse command %s : %s", string(rawCmd), err)
		conn.Write([]byte(err.Error()))
		return
	}

	switch msg.Cmd {
	case CMDSet:
		err = s.handleSetCommand(conn, msg)
	case CMDGet:
		err = s.handleGetCommmand(conn, msg)
	}

	if err != nil {
		log.Printf("server: failed to handle command [%s]: %s", string(rawCmd), err)
		conn.Write([]byte(err.Error()))
		return
	}
}

func (s *Server) handleSetCommand(conn net.Conn, msg *Message) error {
	if err := s.cache.Set(msg.Key, msg.Value, msg.TTL); err != nil {
		return err
	}

	s.sendToFollowers(context.TODO(), msg)

	return nil
}

func (s *Server) handleGetCommmand(conn net.Conn, msg *Message) error {
	val, err := s.cache.Get(msg.Key)
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte(val))

	return err
}

func (s *Server) sendToFollowers(ctx context.Context, msg *Message) error {
	return nil
}
