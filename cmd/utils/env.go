package utils

import "os"

type EnvConfig struct {
	MQTTServer    string
	MQTTPort      string
	MQTTUser      string
	MQTTPassword  string
	EmoncmsApiKey string
	EmoncmsUrl    string
}

func LoadEnvConfig() *EnvConfig {
	return &EnvConfig{
		MQTTServer:    getEnv("MQTT_SERVER", "localhost"),
		MQTTPort:      getEnv("MQTT_PORT", "8883"),
		MQTTUser:      getEnv("MQTT_USER", ""),
		MQTTPassword:  getEnv("MQTT_PASSWORD", ""),
		EmoncmsApiKey: getEnv("EMONCMS_API_KEY", ""),
		EmoncmsUrl:    getEnv("EMONCMS_URL", "https://emoncms.org"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
