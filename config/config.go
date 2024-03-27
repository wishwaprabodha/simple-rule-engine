package config

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type Config struct {
	PhoneNumbers []string `yaml:"PhoneNumbers"`
}

type Instance struct {
	Config Config
	mu     sync.RWMutex
}

var configFile = "config/config.yaml"
var ConfigInstance = &Instance{}

func (c *Instance) LoadConfig() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Read config file
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Unmarshal YAML data into config struct
	err = yaml.Unmarshal(data, &c.Config)
	if err != nil {
		log.Fatalf("Error unmarshalling config data: %v", err)
	}

	log.Println("Config loaded successfully")
}

func (c *Instance) writeConfigToFile() error {
	data, err := yaml.Marshal(&c.Config)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(configFile, data, 0644); err != nil {
		return err
	}

	return nil
}

func UpdateConfigHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	newConfig := Config{}
	if err := yaml.Unmarshal(body, &newConfig); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	ConfigInstance.mu.Lock()

	ConfigInstance.Config.PhoneNumbers = append(ConfigInstance.Config.PhoneNumbers, newConfig.PhoneNumbers...)
	ConfigInstance.mu.Unlock()

	if err := ConfigInstance.writeConfigToFile(); err != nil {
		http.Error(w, "Failed to write updated config to file", http.StatusInternalServerError)
		return
	}

	// Real Time Config Update
	ConfigInstance.LoadConfig()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Config updated successfully"))
}

func GetConfigHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	ConfigInstance.mu.RLock()
	defer ConfigInstance.mu.RUnlock()

	data, err := json.Marshal(&ConfigInstance.Config)
	if err != nil {
		http.Error(w, "Failed to marshal config data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
