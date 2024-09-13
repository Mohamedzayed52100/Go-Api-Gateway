package guest

import (
	"context"
	"net/http"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SGuest) GetReservationByID() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		id := c.Param("reservationId")

		proto, err := s.ReservationClient.GetReservationByID(ctx, &guestProto.GetReservationByIDRequest{
			Id: common.ConvertStringToInt(id),
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildReservationAPIResponse(proto.GetResult()))
	}
}
