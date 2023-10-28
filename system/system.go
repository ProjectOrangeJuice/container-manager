package system

import (
	"fmt"
	"log"
	"time"

	"github.com/ProjectOrangeJuice/vm-manager-server/connection"
)

type System struct {
	clients connection.Clients
}

func RunSystem(containerManager connection.Clients) {
	ticker := time.NewTicker(30 * time.Second)
	for {
		<-ticker.C
		log.Printf("Ticking over the containers")
		clients := containerManager.GetActiveClients()

		for _, client := range clients {
			fmt.Fprint(client.Conn, "STORAGE_INFO\n")
			fmt.Fprint(client.Conn, "SYSTEM_INFO\n")
		}

	}
}
