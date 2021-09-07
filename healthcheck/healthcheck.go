package healthcheck

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "world",
	})
}
