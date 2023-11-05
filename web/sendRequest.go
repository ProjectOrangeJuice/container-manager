package web

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (w webServer) SendUpdateRequest(c *gin.Context) {
	id := c.Param("id")
	log.Printf("Got update request for %s", id)
	err := w.clients.SendUpdateRequest(id)
	if err != nil {
		log.Printf("Error sending update request: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
}
