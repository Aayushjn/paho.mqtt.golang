/*
 * Copyright (c) 2021 IBM Corp and others.
 *
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v2.0
 * and Eclipse Distribution License v1.0 which accompany this distribution.
 *
 * The Eclipse Public License is available at
 *    https://www.eclipse.org/legal/epl-2.0/
 * and the Eclipse Distribution License is available at
 *   http://www.eclipse.org/org/documents/edl-v10.php.
 *
 * Contributors:
 *    Seth Hoenig
 *    Allan Stockdill-Mander
 *    Mike Robertson
 */

package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)

	// cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	// if err != nil {
	// 	panic(err)
	// }

	// tlsConfig := &tls.Config{
	// 	ClientAuth:         tls.NoClientCert,
	// 	ClientCAs:          nil,
	// 	InsecureSkipVerify: true,
	// 	Certificates:       []tls.Certificate{cert},
	// }

	opts := mqtt.NewClientOptions().AddBroker("quic://127.0.0.1:1883").SetClientID("gotrivial").SetTLSConfig(&tls.Config{}).SetUsername("testuser").SetPassword("testpassword").SetStore(mqtt.NewFileStore("/home/aayush/paho.mqtt.golang/cmd/connect_only/messages"))
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	start := time.Now()
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	end := time.Now()
	fmt.Printf("CONN time taken: %v\n", end.Sub(start))
	fmt.Println("connected!!!")

	if token := c.Subscribe("go-mqtt/sample", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	fmt.Println("subscribed!!!")

	for i := 0; i < 5; i++ {
		text := fmt.Sprintf("this is msg #%d!", i)
		token := c.Publish("go-mqtt/sample", 0, false, text)
		token.Wait()
		fmt.Println("published!!!", token.Error())
	}

	time.Sleep(6 * time.Second)

	if token := c.Unsubscribe("go-mqtt/sample"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	fmt.Println("unsubscribed!!!")

	c.Disconnect(250)

	time.Sleep(1 * time.Second)
}
