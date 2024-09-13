package guest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) GetGuestUpcomingReservations() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		guestId := c.Param("guestId")

		proto, err := s.GuestClient.GetGuestUpcomingReservations(ctx, &guestProto.GetGuestUpcomingReservationsRequest{
			GuestId: common.ConvertStringToInt(guestId),
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildReservationsAPIResponse(proto.GetResult()))
	}
}
