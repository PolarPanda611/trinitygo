package runtime

import (
	"fmt"

	"github.com/PolarPanda611/trinitygo/application"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
)

// New runtime middleware
func New(app application.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, v := range app.RuntimeKeys() {
			keyValue := c.GetHeader(v.GetKeyName())
			if keyValue == "" {
				if v.GetRequired() {
					if app.ResponseFactory() != nil {
						c.JSON(400, app.ResponseFactory()(400, map[string]string{
							"code":    codes.Internal.String(),
							"message": fmt.Sprintf("runtime key %v is required ", v.GetKeyName()),
						},
							nil,
						))
					} else {
						c.AbortWithStatusJSON(400, map[string]string{
							"code":    codes.Internal.String(),
							"message": fmt.Sprintf("runtime key %v is required ", v.GetKeyName()),
						})
					}
					return
				}
				c.Set(v.GetKeyName(), v.GetDefaultValue())
				if v.IsLog() {
					c.Header(v.GetKeyName(), v.GetDefaultValue())
				}
			} else {
				c.Set(v.GetKeyName(), keyValue)
				c.Header(v.GetKeyName(), keyValue)
			}
		}
		c.Next()
	}
}
