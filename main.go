package main

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rmhyde/ashp-metrics/cmd/utils"
)

func main() {
	config := utils.LoadEnvConfig()
	broker := fmt.Sprintf("tls://%s:%s", config.MQTTServer, config.MQTTPort)
	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID("ashp-metrics")
	opts.SetClientID("ashp-metrics-client")
	opts.SetUsername(config.MQTTUser)
	opts.SetPassword(config.MQTTPassword)

	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("MQTT connection error:", token.Error())
		return
	}
	subscribe(client)

	time.Sleep(60 * 1e9)
	client.Disconnect(250)
}

func subscribe(client mqtt.Client) {
	// subscribe to the same topic, that was published to, to receive the messages
	topic := "espaltherma/ATTR"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	// Check for errors during subscribe (More on error reporting https://pkg.go.dev/github.com/eclipse/paho.mqtt.golang#readme-error-handling)
	if token.Error() != nil {
		fmt.Printf("Failed to subscribe to topic")
		panic(token.Error())
	}
	fmt.Printf("Subscribed to topic: %s", topic)
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v", err)
}
