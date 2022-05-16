package config

import (
	"encoding/base64"
	"os"

	log "github.com/sirupsen/logrus"
)

type Configuration struct {
	GcpPubSubConfig
	ProtoDescFilePath string
}
type GcpPubSubConfig struct {
	CredentialJson string
}

func ReadConfiguration() (Configuration, error) {
	log.Printf("Initialize configs.")
	configuration := Configuration{}

	//pubsub
	//configuration.GcpPubSubConfig.EmulatorHost = GetEnv("PUBSUB_EMULATOR_HOST", "0.0.0.0:8681")
	configuration.GcpPubSubConfig.CredentialJson = base64Decode(GetEnv("GOOGLE_CREDENTIAL_JSON", "{}"))

	configuration.ProtoDescFilePath = GetEnv("PROTO_DESC_PATH", "/Users/myuser/Downloads/descriptor.desc")
	return configuration, nil
}

func base64Decode(str string) string {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "{}"
	}
	return string(data)
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	log.Printf("Could not load environment var %s, using default %s", key, fallback)
	return fallback
}
