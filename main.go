package main

import (
	"log"
	"net"
	"time"

	"github.com/AshirwadPradhan/ggcache/cache"
)

func main() {
	opts := ServerOpts{
		ListenAddr: ":3000",
		IsLeader:   true,
	}

	go func() {
		time.Sleep(time.Second * 2)
		conn, err := net.Dial("tcp", ":3000")
		if err != nil {
			log.Fatal(err)
		}

		conn.Write([]byte("SET Foo Bar 1000"))
	}()

	server := NewServer(opts, cache.New())
	server.Start()

}
