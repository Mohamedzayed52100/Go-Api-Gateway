package guest

import (
	"context"
	"net/http"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SGuest) GetAllReservationLogs() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		ReservationId := c.Param("reservationId")

		proto, err := s.ReservationLogClient.GetAllReservationLogs(ctx, &guestProto.GetAllReservationLogsRequest{
			ReservationId: common.ConvertStringToInt(ReservationId),
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildAllReservationLogsResponse(proto.GetResult()))
	}
}
