package middlewares

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/xybor/xychat/helpers"
)

// ApplyCORSHeader adds the Access-Control-Allow-Origin to the header of
// response.
func ApplyCORSHeader() gin.HandlerFunc {
	config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		AllowWebSockets:  true,
		MaxAge:           2 * time.Hour,
	}

	domains, err := helpers.ReadEnv("CORS")
	if err != nil {
		config.AllowAllOrigins = true
	} else {
		domainList := strings.Split(domains, ";")

		for i, d := range domainList {
			domainList[i] = strings.Trim(d, " ")
		}

		config.AllowOrigins = domainList
	}

	return cors.New(config)
}
