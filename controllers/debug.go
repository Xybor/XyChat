package controllers

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
)

const (
	GET int = iota
	JSON
	POST
)

var method = GET

func SetReceivingMethod(m int) {
	if m != GET && m != POST && m != JSON {
		log.Fatal("invalid receiving method")
	}

	method = m
}

func getGET(c *gin.Context, key string) (string, error) {
	v, ok := c.GetQuery(key)
	if !ok {
		return "", errors.New("invalid key")
	}

	return v, nil
}

func getPOST(c *gin.Context, key string) (string, error) {
	v, ok := c.GetPostForm(key)
	if !ok {
		return "", errors.New("invalid key")
	}

	return v, nil
}

func getJSON(c *gin.Context, key string) (string, error) {
	var param map[string]string

	if c.ShouldBindJSON(&param) != nil {
		return "", errors.New("invalid body")
	}

	v, ok := param[key]
	if !ok {
		return "", errors.New("invalid key")
	}

	return v, nil
}

func GetParam(c *gin.Context, key string) (string, error) {
	switch method {
	case GET:
		return getGET(c, key)
	case POST:
		return getPOST(c, key)
	case JSON:
		return getJSON(c, key)
	default:
		return getGET(c, key)
	}
}
