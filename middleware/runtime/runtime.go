package runtime

import (
	"fmt"

	"github.com/PolarPanda611/trinitygo/application"
	"github.com/gin-gonic/gin"
)

// New runtime middleware
func New(app application.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, v := range app.RuntimeKeys() {
			fmt.Println(v)
			keyValue := c.GetString(v.GetKeyName())
			if v.GetRequired() {
				if keyValue == "" {
				}

			}

		}

		c.Next()
	}
}
