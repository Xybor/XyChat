package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/xybor/xychat/helpers"
)

// ApplyAPIHeader adds some headers for API response.
func ApplyAPIHeader(c *gin.Context) {
	c.Writer.Header().Add("Content-Type", "application/json")
}

// ApplyCORSHeader adds the Access-Control-Allow-Origin to the header of
// response.
func ApplyCORSHeader(c *gin.Context) {
	cors_url := fmt.Sprintf("%s://%s%s", helpers.MustReadEnv("SCHEMA"), helpers.MustReadEnv("DOMAIN"), helpers.MustReadEnv("CLIENT_PORT"))
	c.Writer.Header().Add("Access-Control-Allow-Origin", cors_url)
	c.Writer.Header().Add("Access-Control-Allow-Credentials", "true")
}
