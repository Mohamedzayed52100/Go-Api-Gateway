package guest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) DeleteReservationVisitors() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		reservationId := common.ConvertStringToInt(c.Param("reservationId"))
		var input openapi.DeleteReservationVisitorsRequest
		common.ResolveRequestBody(c, &input)

		proto, err := s.GuestClient.DeleteReservationVisitors(ctx, &guestProto.DeleteReservationVisitorsRequest{
			ReservationId: reservationId,
			VisitorIds:    input.Visitors,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, openapi.DeleteReservationVisitorsResponse{
			Code:    proto.GetCode(),
			Message: proto.GetMessage(),
		})
	}
}
