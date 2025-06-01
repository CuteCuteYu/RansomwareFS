package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/key", generateKeysHandler)
	r.GET("/get_key", getKeysHandler)
	r.Run(":8080")
}
