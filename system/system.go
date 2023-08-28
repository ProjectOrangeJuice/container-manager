package system

import (
	"fmt"
	"log"
	"time"

	"github.com/ProjectOrangeJuice/vm-manager-server/vm"
)

type System struct {
	VMs vm.Container
}

func RunSystem(containerManager vm.Container) {
	ticker := time.NewTicker(30 * time.Second)
	for {
		<-ticker.C
		log.Printf("Ticking over the containers")
		clients := containerManager.GetAllClients()

		for _, client := range clients {
			fmt.Fprint(client.Conn, "STORAGE_INFO\n")
			fmt.Fprint(client.Conn, "SYSTEM_INFO\n")
		}

	}
}
