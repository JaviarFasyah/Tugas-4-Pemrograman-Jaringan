package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("[ %s ]\n", msg.Topic())
	fmt.Printf("%s\n", msg.Payload())
}

func main() {
	// Config
	opts := MQTT.NewClientOptions().AddBroker("localhost:1883")
	opts.SetClientID("subscriber")
	opts.SetDefaultPublishHandler(f)

	// Connect ke broker
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Subscribe ke topik(tpc) yang diinginkan
	tpc := flag.String("topic", "#", "the topic")
	flag.Parse()
	if *tpc == "" {
		*tpc = "#"
	}

	// Subscribe ke topik tpc
	if token := c.Subscribe(*tpc, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	// Subscribe selama 1 menit
	time.Sleep(60 * time.Second)

	// Unsubscribe
	if token := c.Unsubscribe(*tpc); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	c.Disconnect(250)
}
