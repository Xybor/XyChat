package helpers

import "github.com/gin-gonic/gin"

// SetCookie is a shortcut of ctx.SetCookie with path, domain, secure, and
// httponly parameters are automatically filled.
func SetCookie(ctx *gin.Context, name, value string, maxage int) {
	ctx.SetCookie(
		name,
		value,
		maxage,
		"/",
		MustReadEnv("DOMAIN"),
		false,
		false,
	)
}
