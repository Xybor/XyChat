package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StaticUIHandler(ctx *gin.Context) {
	// location := url.URL{Path: "/"}
	// ctx.Redirect(http.StatusAccepted, location.RequestURI())
	ctx.HTML(http.StatusAccepted, "index.html", gin.H{})
}