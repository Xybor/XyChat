package api

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Data *interface{}            `json:"data,omitempty"`
	Meta *map[string]interface{} `json:"meta,omitempty"`
}

func CreateAPIResponse(
	data *interface{},
	meta *map[string]interface{},
) APIResponse {
	return APIResponse{Data: data, Meta: meta}
}

func ApplyAPIHeader(c *gin.Context) {
	c.Writer.Header().Add("Content-Type", "application/json")
}
