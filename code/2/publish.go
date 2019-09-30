package main

import (
	"flag"
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("[ %s ]\n", msg.Topic())
	fmt.Printf("%s\n", msg.Payload())
}

func main() {
	// Config
	opts := MQTT.NewClientOptions().AddBroker("localhost:1883")
	opts.SetClientID("publisher")
	opts.SetDefaultPublishHandler(f)

	// Connect ke broker
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Topik(tpc) & pesan(msg) yang diinginkan
	tpc := flag.String("topic", "foo", "the topic")
	msg := flag.String("message", "", "the message")
	flag.Parse()
	if *tpc == "" {
		*tpc = "foo"
	}

	// Publish pesan ke tpc
	token := c.Publish(*tpc, 0, false, *msg)
	token.Wait()

	c.Disconnect(250)
}
