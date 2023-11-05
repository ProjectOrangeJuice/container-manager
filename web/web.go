package web

import (
	"log"
	"net/http"

	"github.com/ProjectOrangeJuice/vm-manager-server/connection"
	"github.com/gin-contrib/static"
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
	// deal with cors and preflight requests
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		}
		c.Next()
	})

	// Serve static files
	r.Use(static.Serve("/", static.LocalFile("./static", false)))

	r.GET("/api/list", w.handleListAPI)
	r.POST("/api/waiting/:id", w.handleWaitingAPI)

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

type waitingResponse struct {
	Allow bool
}

func (w webServer) handleWaitingAPI(c *gin.Context) {

	id := c.Param("id")
	var a waitingResponse
	c.BindJSON(&a)

	log.Printf("Waiting client %s will be allowed? %v", id, a.Allow)

	w.clients.DealWithWaiting(id, a.Allow)
}

func (w webServer) handleListAPI(c *gin.Context) {

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
