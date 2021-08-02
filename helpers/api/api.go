package api

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Data *interface{}            `json:"data,omitempty"`
	Meta *map[string]interface{} `json:"meta,omitempty"`
}

func ApplyAPIHeader(c *gin.Context) {
	c.Writer.Header().Add("Content-Type", "application/json")
}
