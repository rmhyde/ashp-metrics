package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rmhyde/ashp-metrics/cmd/utils"
	"github.com/rmhyde/ashp-metrics/internal/emoncms"
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

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Looping... Press Ctrl+C to break.")
	done := false
	for !done {
		select {
		case <-sigs:
			fmt.Println("\nBreak signal received!")
			done = true
		default:
			time.Sleep(60 * 1e9)
		}
	}

	client.Disconnect(250)
	fmt.Println("Cleaned up and exited.")
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
	fmt.Printf("Subscribed to topic: %s\n", topic)
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	var payload EspAlthermaPayload
	err := json.Unmarshal(msg.Payload(), &payload)
	if err != nil {
		fmt.Printf("Error unmarshalling payload: %v\n", err)
		return
	}

	convertedPayload := ConvertToEmoncmsPayload(payload)
	fmt.Printf("Parsed payload: %+v\n", convertedPayload)

	config := utils.LoadEnvConfig()
	c := emoncms.NewClient(config.EmoncmsUrl, config.EmoncmsApiKey)
	err = emoncms.UpdateInput(c, "espaltherma", convertedPayload)
	if err != nil {
		fmt.Printf("Error updating input: %v\n", err)
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v", err)
}

func ConvertToEmoncmsPayload(payload EspAlthermaPayload) EmoncmsPayload {
	return EmoncmsPayload{
		OutdoorAirTemp:            payload.OutdoorAirTemp,
		INVPrimaryCurrent:         payload.INVPrimaryCurrent,
		INVSecondaryCurrent:       payload.INVSecondaryCurrent,
		IUOperationMode:           payload.IUOperationMode,
		DHWSetpoint:               payload.DHWSetpoint,
		DHWTankTemp:               payload.DHWTankTemp,
		LWSetpointMain:            payload.LWSetpointMain,
		LeavingWaterTempBeforeBUH: payload.LeavingWaterTempBeforeBUH,
		LeavingWaterTempAfterBUH:  payload.LeavingWaterTempAfterBUH,
		InletWaterTemp:            payload.InletWaterTemp,
		IndoorAmbientTemp:         payload.IndoorAmbientTemp,
		RTSetpoint:                payload.RTSetpoint,
		FlowRate:                  payload.FlowRate,
		Pressure:                  payload.Pressure,
	}
}

type EmoncmsPayload struct {
	OutdoorAirTemp            float64 `json:"OutdoorTemp"`
	INVPrimaryCurrent         float64 `json:"InvPrimaryCurrent"`
	INVSecondaryCurrent       float64 `json:"InvSecondaryCurrent"`
	IUOperationMode           string  `json:"IUOperationMode"`
	DHWSetpoint               float64 `json:"DHWSetpoint"`
	DHWTankTemp               float64 `json:"DHWTemp"`
	LWSetpointMain            float64 `json:"LWSetpointMain"`
	LeavingWaterTempBeforeBUH float64 `json:"LWTempbfBUH"`
	LeavingWaterTempAfterBUH  float64 `json:"LWTempAfBUH"`
	InletWaterTemp            float64 `json:"RWTemp"`
	IndoorAmbientTemp         float64 `json:"RoomTemp"`
	RTSetpoint                float64 `json:"RoomSetpoint"`
	FlowRate                  float64 `json:"FlowRate"`
	Pressure                  float64 `json:"Pressure"`
}

type EspAlthermaPayload struct {
	OperationMode             string  `json:"Operation Mode"`
	ThermostatONOFF           string  `json:"Thermostat ON/OFF"`
	ErrorType                 string  `json:"Error type"`
	ErrorCode                 string  `json:"Error Code"`
	OutdoorAirTemp            float64 `json:"R1T-Outdoor air temp."`
	INVPrimaryCurrent         float64 `json:"INV primary current (A)"`
	INVSecondaryCurrent       float64 `json:"INV secondary current (A)"`
	IUOperationMode           string  `json:"I/U operation mode"`
	DHWSetpoint               float64 `json:"DHW setpoint"`
	LWSetpointMain            float64 `json:"LW setpoint (main)"`
	LeavingWaterTempBeforeBUH float64 `json:"Leaving water temp. before BUH (R1T)"`
	LeavingWaterTempAfterBUH  float64 `json:"Leaving water temp. after BUH (R2T)"`
	InletWaterTemp            float64 `json:"Inlet water temp.(R4T)"`
	DHWTankTemp               float64 `json:"DHW tank temp. (R5T)"`
	IndoorAmbientTemp         float64 `json:"Indoor ambient temp. (R1T)"`
	RTSetpoint                float64 `json:"RT setpoint"`
	FlowRate                  float64 `json:"Flow sensor (l/min)"`
	Pressure                  float64 `json:"Water pressure"`
}
