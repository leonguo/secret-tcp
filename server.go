package main

import (
	"test/tcp3/connect"
)

func main() {
	conf := connect.Conf{
		Address:      "127.0.0.1:9999",
		MaxConnCount: 100,
		AcceptCount:  1,
	}

	server := connect.NewTCPServer(conf)
	server.Start()
}
