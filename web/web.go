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
	r.LoadHTMLGlob("web/templates/*")
	r.GET("/", w.handleList)
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

func (w webServer) handleList(c *gin.Context) {
	clients := w.clients.GetAllClients()
	webResult := []clientResult{}
	for _, cc := range clients {
		client := clientResult{
			Name:    cc.Name,
			CPU:     cc.System.CPUUseage,
			Memory:  (float64(cc.System.TotalMemory-cc.System.FreeMemory) / float64(cc.System.TotalMemory)) * 100,
			Storage: []storageResult{},
		}
		for _, storage := range cc.Storage {
			client.Storage = append(client.Storage, storageResult{
				Name:      storage.Name,
				Mount:     storage.Mount,
				SpaceUsed: (storage.UsedSpace / storage.TotalSpace) * 100,
			})
		}
		webResult = append(webResult, client)
	}

	log.Printf("Web result: %+v", webResult)

	c.HTML(http.StatusOK, "index", webResult)
}

func (w webServer) handleListAPI(c *gin.Context) {
	// Set CORS headers
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type")

	clients := w.clients.GetAllClients()
	webResult := []clientResult{}
	for _, cc := range clients {
		client := clientResult{
			Name:    cc.Name,
			CPU:     cc.System.CPUUseage,
			Memory:  (float64(cc.System.TotalMemory-cc.System.FreeMemory) / float64(cc.System.TotalMemory)) * 100,
			Storage: []storageResult{},
		}
		for _, storage := range cc.Storage {
			client.Storage = append(client.Storage, storageResult{
				Name:      storage.Name,
				Mount:     storage.Mount,
				SpaceUsed: (storage.UsedSpace / storage.TotalSpace) * 100,
			})
		}
		webResult = append(webResult, client)
	}

	log.Printf("Web result: %+v", webResult)

	c.JSON(http.StatusOK, webResult)
}
