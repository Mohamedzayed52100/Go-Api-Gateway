package general

import (
	"context"
	"net/http"

	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"

	"github.com/gin-gonic/gin"
)

func (s *SGeneral) Index() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	}
}
