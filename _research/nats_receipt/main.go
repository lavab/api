package main

import (
	"github.com/apcera/nats"
	"log"
)

func main() {
	nc, err := nats.Connect("tcp://127.0.0.1:4222")
	if err != nil {
		log.Fatal(err)
	}

	c, err := nats.NewEncodedConn(nc, "json")
	if err != nil {
		log.Fatal(err)
	}

	c.Publish("receipt", struct {
		ID    string `json:"id"`
		Owner string `json:"owner"`
	}{
		ID:    "helloworld",
		Owner: "k5TuCXomSMEeCdXw1aXl",
	})

	log.Print("x")
}
