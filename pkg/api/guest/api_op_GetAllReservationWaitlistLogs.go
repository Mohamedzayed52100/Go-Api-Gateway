package guest

import (
	"context"
	"net/http"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SGuest) GetAllReservationWaitlistLogs() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		reservationWaitlistId := c.Param("reservationWaitlistId")

		proto, err := s.ReservationLogClient.GetAllReservationWaitlistLogs(ctx, &guestProto.GetAllReservationWaitlistLogsRequest{
			ReservationWaitlistId: common.ConvertStringToInt(reservationWaitlistId),
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildAllReservationWaitlistLogsResponse(proto.GetResult()))
	}
}
