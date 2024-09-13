package settings

import (
	"context"
	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
	"net/http"
)

func (s *SSettings) CreateSeatingArea() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.CreateSeatingAreaRequest
		common.ResolveRequestBody(c, &input)

		proto, err := s.SeatingAreaClient.CreateSeatingArea(ctx, &settingsProto.CreateSeatingAreaRequest{
			Params: &settingsProto.SeatingAreaParams{
				Name: input.Name,
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildSeatingAreaResponse(proto.GetResult()))
	}
}
