package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	healthy := false
	sigChan := make(chan os.Signal, 1)
	doneChan := make(chan bool)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, world!",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		if healthy {
			c.String(200, "healthy")
		} else {
			c.String(500, "unhealthy")
		}
	})

	r.GET("/:name", func(c *gin.Context) {
		name := c.Param("name")

		c.JSON(200, gin.H{
			"message": fmt.Sprintf("Hello, %s!", name),
		})
	})

	go func() {
		r.Run(fmt.Sprintf(":%s", port))
	}()

	go func() {
		sig := <-sigChan // This blocks until a signal is received.
		fmt.Printf("Received signal: %s, shutting down...\n", sig)
		// Perform any cleanup here

		// Signal that cleanup is done and the program can exit.
		doneChan <- true
	}()

	healthy = true
	<-doneChan
	healthy = false
}
