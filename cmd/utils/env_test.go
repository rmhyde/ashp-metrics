package utils

import (
	"os"
	"testing"
)

func TestLoadEnvConfig(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected *EnvConfig
	}{
		{
			name:    "default values",
			envVars: map[string]string{},
			expected: &EnvConfig{
				MQTTServer:   "localhost",
				MQTTPort:     "8883",
				MQTTUser:     "",
				MQTTPassword: "",
			},
		},
		{
			name: "custom values",
			envVars: map[string]string{
				"MQTT_SERVER":   "broker.example.com",
				"MQTT_PORT":     "1883",
				"MQTT_USER":     "testuser",
				"MQTT_PASSWORD": "testpass",
			},
			expected: &EnvConfig{
				MQTTServer:   "broker.example.com",
				MQTTPort:     "1883",
				MQTTUser:     "testuser",
				MQTTPassword: "testpass",
			},
		},
		{
			name: "partial override",
			envVars: map[string]string{
				"MQTT_SERVER": "custom.broker",
			},
			expected: &EnvConfig{
				MQTTServer:   "custom.broker",
				MQTTPort:     "8883",
				MQTTUser:     "",
				MQTTPassword: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key := range map[string]bool{"MQTT_SERVER": true, "MQTT_PORT": true, "MQTT_USER": true, "MQTT_PASSWORD": true} {
				os.Unsetenv(key)
			}

			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			got := LoadEnvConfig()

			if got.MQTTServer != tt.expected.MQTTServer || got.MQTTPort != tt.expected.MQTTPort ||
				got.MQTTUser != tt.expected.MQTTUser || got.MQTTPassword != tt.expected.MQTTPassword {
				t.Errorf("LoadEnvConfig() = %+v, want %+v", got, tt.expected)
			}

			for key := range tt.envVars {
				os.Unsetenv(key)
			}
		})
	}
}
