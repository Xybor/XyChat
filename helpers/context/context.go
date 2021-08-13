package context

import (
	"errors"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xybor/xychat/helpers"
)

const (
	GET  = "GET"
	POST = "POST"
)

var method = GET

// SetRetrievingMethod sets the method for get query parameters in the http
// request.  The method should be GET or POST.  The affected functions
// include: RetrieveQuery-(), MustRetrieveQuery().
//
// If the method is GET, the parameter should be get by ctx.GetQuery().
//
// If the method is POST, the parameter should be get by ctx.GetPostForm().
func SetRetrievingMethod(m string) {
	if m != GET && m != POST {
		log.Fatalln("invalid receiving method: " + m)
	}

	method = m
}

func retrieveGET(ctx *gin.Context, key string) (string, error) {
	v, ok := ctx.GetQuery(key)
	if !ok {
		return "", errors.New("invalid get key")
	}

	return v, nil
}

func retrievePOST(ctx *gin.Context, key string) (string, error) {
	v, ok := ctx.GetPostForm(key)
	if !ok {
		return "", errors.New("invalid post key")
	}

	return v, nil
}

// RetrieveQuery bases on the method set by SetRetrievingMethod to get the
// query parameters in the suitable place.
func RetrieveQuery(ctx *gin.Context, key string) (string, error) {
	switch method {
	case GET:
		return retrieveGET(ctx, key)
	case POST:
		return retrievePOST(ctx, key)
	default:
		return retrieveGET(ctx, key)
	}
}

// RetrieveQueryAsPUint retrieves the query parameter and converts it to *uint
// (from string).  If key is invalid, the return value is nil and there is no
// error.
func RetrieveQueryAsPUint(ctx *gin.Context, key string) (*uint, error) {
	svalue, err := RetrieveQuery(ctx, key)
	if err != nil {
		return nil, nil
	}

	value64, err := strconv.ParseUint(svalue, 10, 64)
	if err != nil {
		return nil, err
	}

	value := uint(value64)
	return &value, nil
}

// RetrieveQueryAsPString retrieves the query parameter and converts it to
// *string. If key is invalid, the return value is nil.
func RetrieveQueryAsPString(ctx *gin.Context, key string) *string {
	value, err := RetrieveQuery(ctx, key)
	if err != nil {
		return nil
	}

	return &value
}

// MustRetrieveQuery is equipvalent to RetrieveQuery and log.Panic if a error
// occurs.
func MustRetrieveQuery(ctx *gin.Context, key string) string {
	value, err := RetrieveQuery(ctx, key)

	if err != nil {
		log.Panicln(err)
	}

	return value
}

// GetURLParamAsUint gets the param in URL by using ctx.Params.Get and converts
// it to uint (from string).
func GetURLParamAsUint(ctx *gin.Context, key string) (uint, error) {
	value, ok := ctx.Params.Get(key)
	if !ok {
		return 0, errors.New("invalid key")
	}

	i64, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, errors.New("param can't be parsed as uint")
	}

	return uint(i64), nil
}

// GetUID gets UID from the context by ctx.Get() and convert it to *uint.  If
// an error is occurs, nil value will be returned.
func GetUID(ctx *gin.Context) *uint {
	tmp, ok := ctx.Get("UID")
	if !ok {
		return nil
	}

	id := tmp.(uint)

	return &id
}

// SetCookie is a shortcut of ctx.SetCookie with path, domain, secure, and
// httponly parameters are automatically filled.
func SetCookie(ctx *gin.Context, name, value string, maxage int) {
	ctx.SetCookie(
		name,
		value,
		maxage,
		"/",
		helpers.MustReadEnv("DOMAIN"),
		false,
		false,
	)
}
