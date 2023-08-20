package web

import (
	cell "container-manager/server/container"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type webServer struct {
	containerManager cell.Container
}

func StartWebServer(containerManager cell.Container) {
	w := webServer{
		containerManager: containerManager,
	}

	r := gin.Default()
	r.LoadHTMLGlob("web/templates/*")
	r.GET("/", w.handleList)

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
	clients := w.containerManager.GetAllClients()
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
