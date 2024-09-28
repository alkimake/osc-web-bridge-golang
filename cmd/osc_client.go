package main

import (
	"time"

	"math/rand"

	"github.com/hypebeast/go-osc/osc"
)

func main() {
	client := osc.NewClient("239.255.0.1", 9000)

	for {
		msg := osc.NewMessage("/osc/address")
		msg.Append(int32(rand.Intn(100))) // Random integer
		msg.Append(rand.Intn(2) == 1)     // Random boolean
		msg.Append("hello")               // Static string or can be randomized

		client.Send(msg)
		time.Sleep(3 * time.Second) // Wait for 3 seconds
	}
}
