package middlewares

import "github.com/gin-gonic/gin"

// ApplyAPIHeader adds some headers for API response.
func ApplyAPIHeader(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
}

// ApplyCORSHeader adds the Access-Control-Allow-Origin to the header of
// response.
func ApplyCORSHeader(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
}
