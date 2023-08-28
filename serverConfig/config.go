package serverConfig

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	KeyLocation        string
	ServerAddress      []string
	ClientFingerprints []Fingerprint
}

type Fingerprint struct {
	SerialNumber string
	Name         string
	Fingerprint  string
	AllowConnect bool
}

func FirstRun() error {
	config := Config{
		KeyLocation:   "./keys/",
		ServerAddress: []string{"localhost:8080"},
	}

	// write json to file
	file, err := os.Create("config.json")
	if err != nil {
		return fmt.Errorf("could not create config file, %s", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(config)
	if err != nil {
		return fmt.Errorf("could not encode config file, %s", err)
	}
	return nil
}

func ReadConfig() (Config, bool, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return Config{}, false, fmt.Errorf("could not open config file, %s", err)
	}
	defer file.Close()

	config := Config{}
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return Config{}, true, fmt.Errorf("could not decode config file, %s", err)
	}
	return config, true, nil
}

// Read config and add a new fingerprint
func AddFingerprint(fingerprint Fingerprint) error {
	config, _, err := ReadConfig()
	if err != nil {
		return fmt.Errorf("could not read config, %s", err)
	}
	config.ClientFingerprints = append(config.ClientFingerprints, fingerprint)
	err = WriteConfig(config)
	if err != nil {
		return fmt.Errorf("could not write config, %s", err)
	}
	return nil
}

func WriteConfig(config Config) error {
	file, err := os.Create("config.json")
	if err != nil {
		return fmt.Errorf("could not create config file, %s", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(config)
	if err != nil {
		return fmt.Errorf("could not encode config file, %s", err)
	}
	return nil
}
