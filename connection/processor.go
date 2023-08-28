package connection

import (
	"encoding/json"
	"log"

	"github.com/ProjectOrangeJuice/vm-manager-server/shared"
)

func (c *Client) processEvent(line string) {
	evt := shared.EventData{}
	err := json.Unmarshal([]byte(line), &evt)
	if err != nil {
		log.Printf("Could not process event: %s (%s)", err, line)
		return
	}

	switch evt.Request {
	case "STORAGE":
		c.processStorageEvent(evt.Result)
	case "SYSTEM":
		c.processSystemEvent(evt.Result)
	}
}

type StorageResult struct {
	Name       string
	TotalSpace float64
	UsedSpace  float64
	Mount      string
}

func (c *Client) processSystemEvent(result []byte) {
	log.Printf("Processing system event")
	data := shared.SystemResult{}
	err := json.Unmarshal(result, &data)
	if err != nil {
		log.Printf("Could not process data: %s (%s)", err, result)
		return
	}

	c.System = data
	log.Printf("Client %s has updated their system [%+v]", c.Name, c.System)
}

func (c *Client) processStorageEvent(result []byte) {
	log.Printf("Processing storage event")
	data := []shared.StorageResult{}
	err := json.Unmarshal(result, &data)
	if err != nil {
		log.Printf("Could not process data: %s (%s)", err, result)
		return
	}

	c.Storage = data
	log.Printf("Client %s has updated their storage [%+v]", c.Name, c.Storage)
}
