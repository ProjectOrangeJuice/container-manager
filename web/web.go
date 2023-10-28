package web

import (
	"log"
	"net/http"

	"github.com/ProjectOrangeJuice/vm-manager-server/connection"
	"github.com/gin-gonic/gin"
)

type webServer struct {
	clients connection.Clients
}

func StartWebServer(clients connection.Clients) {
	w := webServer{
		clients: clients,
	}

	r := gin.Default()
	r.OPTIONS("/api/list", w.handleListAPI)
	r.GET("/api/list", w.handleListAPI)

	r.Run(":8081")
}

type clientResult struct {
	Name    string
	CPU     float64
	Memory  float64
	Storage []storageResult
}

type storageResult struct {
	Name      string
	Mount     string
	SpaceUsed float64
}

type listAPIResult struct {
	ActiveClients       []clientResult
	WaitingClients      []connection.ClientDetails
	DisconnectedClients []connection.ClientDetails
}

func (w webServer) handleListAPI(c *gin.Context) {
	// Set CORS headers
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type")

	activeClients := w.clients.GetActiveClients()
	waitingClients := w.clients.GetWaitingClients()
	acceptedClients := w.clients.GetAcceptedClients()

	disconnectedClients := createDisconnectedList(activeClients, acceptedClients)

	listAPIResult := listAPIResult{
		DisconnectedClients: disconnectedClients,
		WaitingClients:      waitingClients,
	}
	for _, cc := range activeClients {
		clientDetail := createClientDetail(*cc)
		listAPIResult.ActiveClients = append(listAPIResult.ActiveClients, clientDetail)
	}

	log.Printf("api result: %+v", listAPIResult)

	c.JSON(http.StatusOK, listAPIResult)
}

func createClientDetail(c connection.Client) clientResult {
	client := clientResult{
		Name:    c.Name,
		CPU:     c.System.CPUUseage,
		Memory:  (float64(c.System.TotalMemory-c.System.FreeMemory) / float64(c.System.TotalMemory)) * 100,
		Storage: []storageResult{},
	}
	for _, storage := range c.Storage {
		client.Storage = append(client.Storage, storageResult{
			Name:      storage.Name,
			Mount:     storage.Mount,
			SpaceUsed: (storage.UsedSpace / storage.TotalSpace) * 100,
		})
	}
	return client
}

func createDisconnectedList(activeClients []*connection.Client, acceptedClients []connection.ClientDetails) []connection.ClientDetails {
	disconnectedClients := []connection.ClientDetails{}
	// if the client is in the accepted list but not in the active list, then it is disconnected
	for _, ac := range acceptedClients {
		found := false
		for _, cc := range activeClients {
			if ac.Name == cc.Name {
				found = true
				break
			}
		}
		if !found {
			disconnectedClients = append(disconnectedClients, ac)
		}
	}
	return disconnectedClients
}
