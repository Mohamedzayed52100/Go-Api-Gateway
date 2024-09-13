package guest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) GetAllWaitingReservations() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		shiftId := c.Query("shiftId")
		date := c.Query("date")

		proto, err := s.ReservationWaitlistClient.GetAllWaitingReservations(ctx, &guestProto.GetWaitingReservationRequest{
			ShiftId: common.ConvertStringToInt(shiftId),
			Date:    date,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildAllReservationWaitlistsResponse(proto.GetResult()))
	}
}
