package settings

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
)

func (s *SSettings) UpdateReservationStatus() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		id := c.Param("statusId")
		var input openapi.UpdateReservationStatusRequest
		common.ResolveRequestBody(c, &input)

		proto, err := s.ReservationStatusClient.UpdateReservationStatus(ctx, &settingsProto.UpdateReservationStatusRequest{
			Params: &settingsProto.ReservationStatusParams{
				Id:       common.ConvertStringToInt(id),
				Name:     input.Name,
				Category: input.Category,
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
