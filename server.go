package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/key", generateKeysHandler)
	r.GET("/get_key", getKeysHandler)

	// Run HTTPS server with SSL certificates
	err := r.RunTLS(":443", "ssl_cert/certificate.crt", "ssl_cert/private.key")
	if err != nil {
		panic("Failed to start HTTPS server: " + err.Error())
	}

	// Optional: Also run HTTP server that redirects to HTTPS
	go func() {
		httpRouter := gin.Default()
		httpRouter.GET("/*path", func(c *gin.Context) {
			c.Redirect(301, "https://"+c.Request.Host+c.Request.RequestURI)
		})
		httpRouter.Run(":80")
	}()
}
