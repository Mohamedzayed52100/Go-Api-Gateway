package guest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) GetReservationOrderByReservationId() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		reservationId := c.Param("reservationId")

		proto, err := s.ReservationClient.GetReservationOrderByReservationID(ctx, &guestProto.GetReservationOrderByReservationIDRequest{
			ReservationId: common.ConvertStringToInt(reservationId),
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildReservationOrderResponse(proto.GetResult()))
	}
}
