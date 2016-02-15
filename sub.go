package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
)

func main() {

	addr := flag.String("addr", ":6379", "redis address")
	format := flag.String("fmt", `{"%s":"%s"}`, "data format")

	flag.Parse()
	args := flag.Args()
	conn, err := redis.Dial("tcp", *addr)
	if err != nil {
		log.Fatal(err)
		return
	}
	psc := redis.PubSubConn{conn}

	for _, channel := range args {
		psc.Subscribe(channel)
	}

	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			fmt.Println(fmt.Sprintf(*format, v.Channel, v.Data))
		case redis.Subscription:
			fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			log.Println(v)
		}
	}
}
