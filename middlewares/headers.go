package middlewares

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/xybor/xychat/helpers"
)

// ApplyAPIHeader adds some headers for API response.
func ApplyAPIHeader(c *gin.Context) {
	c.Writer.Header().Add("Content-Type", "application/json")
}

// ApplyCORSHeader adds the Access-Control-Allow-Origin to the header of
// response.
func ApplyCORSHeader() gin.HandlerFunc {
	domains, err := helpers.ReadEnv("CORS")
	config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		AllowWebSockets:  true,
		MaxAge:           2 * time.Hour,
	}

	if err != nil {
		config.AllowAllOrigins = true
	} else {
		domainList := strings.Split(domains, ";")

		// Strip spaces
		for i, d := range domainList {
			domainList[i] = strings.Trim(d, " ")
		}

		config.AllowOrigins = domainList
	}

	return cors.New(config)
}
