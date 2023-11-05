package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// mqtt.DEBUG = log.New(os.Stdout, "DEBUG", 0)
	mqtt.WARN = log.New(os.Stdout, "WARN ", 0)
	mqtt.ERROR = log.New(os.Stderr, "ERROR ", 0)
	mqtt.CRITICAL = log.New(os.Stderr, "CRITICAL ", 0)

	broker := flag.String("broker", "quic://127.0.0.1:1883", "Broker address (protocol must be QUIC)")
	topic := flag.String("topic", "quic-mqtt/test", "Topic to publish/subscribe to")
	flag.Parse()

	// store := mqtt.NewMemoryStore()
	opts := mqtt.NewClientOptions().
		AddBroker(*broker).
		SetClientID("gotrivial").
		// SetProtocolVersion(4).
		SetKeepAlive(2 * time.Second).
		SetPingTimeout(1 * time.Second).
		// SetStore(store).
		SetTLSConfig(&tls.Config{}).
		SetUsername("testuser").
		SetPassword("testpassword").
		SetDefaultPublishHandler(func(c mqtt.Client, m mqtt.Message) {
			mqtt.WARN.Printf("Received: %s from topic: %s\n", m.Payload(), m.Topic())
		})

	c := mqtt.NewClient(opts)
	start := time.Now()
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	end := time.Now()
	fmt.Printf("CONN time taken: %v\n", end.Sub(start))

	if token := c.Subscribe(*topic, 0, nil); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	token := c.Publish(*topic, 0, false, "Hello MQTT!")
	token.Wait()

	time.Sleep(1 * time.Second)

	if token := c.Unsubscribe(*topic); token.Wait() && token.Error() != nil {
		mqtt.CRITICAL.Println(token.Error())
		os.Exit(1)
	}

	time.Sleep(1 * time.Second)
	c.Disconnect(250)

	c = mqtt.NewClient(opts)
	start = time.Now()
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		mqtt.CRITICAL.Println(token.Error())
		os.Exit(1)
	}
	end = time.Now()
	fmt.Printf("CONN time taken: %v\n", end.Sub(start))
	time.Sleep(1 * time.Second)
	c.Disconnect(250)
	// fmt.Printf("Time taken: %v\n", end.Sub(start))
}
