package settings

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
)

func (s *SSettings) CreateReservationStatus() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.CreateReservationStatusRequest
		common.ResolveRequestBody(c, &input)

		proto, err := s.ReservationStatusClient.CreateReservationStatus(ctx, &settingsProto.CreateReservationStatusRequest{
			Params: &settingsProto.ReservationStatusParams{
				Name:     input.Name,
				Category: input.Category,
				Type:     input.Type,
				Color:    input.Color,
				Icon:     input.Icon,
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, buildReservationStatusResponse(proto.GetResult()))
	}
}
