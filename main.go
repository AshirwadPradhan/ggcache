package main

import "github.com/AshirwadPradhan/ggcache/cache"

func main() {
	opts := ServerOpts{
		ListenAddr: ":3000",
		IsLeader:   true,
	}
	server := NewServer(opts, cache.New())
	server.Start()

}
