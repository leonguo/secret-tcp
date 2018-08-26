package main

import (
	"net"
	"log"
	"time"
	"test/tcp3/connect"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:6878")
	if err != nil {
		log.Println("Error dialing", err.Error())
		return
	}

	codec := connect.NewCodec(conn)

	codec.Encode(connect.Message{1, []byte("{\"arg\":\"GetNodeInfo\"}")}, 2*time.Second)
	//codec.Encode(connect.Message{1, []byte("ok fuck")}, 2*time.Second)

	_, err = codec.Read()
	if err != nil {
		log.Println(err)
	}

	for {
		message, ok := codec.Decode()
		if ok {
			spew.Dump(message)
			continue
		}
		break
	}
}