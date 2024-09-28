package main

import "github.com/hypebeast/go-osc/osc"

func main() {
	client := osc.NewClient("239.255.0.1", 9000)
	msg := osc.NewMessage("/osc/address")
	msg.Append(int32(111))
	msg.Append(true)
	msg.Append("hello")
	client.Send(msg)
}
