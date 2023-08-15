package system

import (
	cell "container-manager/server/container"
	"fmt"
	"log"
	"time"
)

type System struct {
	CM cell.Container
}

func RunSystem(containerManager cell.Container) {
	ticker := time.NewTicker(10 * time.Second)
	for {
		<-ticker.C
		log.Printf("Ticking over the containers")
		clients := containerManager.GetAllClients()

		for _, client := range clients {
			fmt.Fprint(client.Conn, "Tick\n")
		}

	}
}
