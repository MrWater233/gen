package gen

import (
	"log"
	"time"
)

// Middleware for calculate resolution time
func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		// execuate router handler
		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
