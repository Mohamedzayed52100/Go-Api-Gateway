package guest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) AddWaitingReservationNote() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		reservationWaitlistId := common.ConvertStringToInt(c.Param("reservationWaitlistId"))

		var input openapi.ReservationWaitlistNote
		common.ResolveRequestBody(c, &input)

		proto, err := s.ReservationWaitlistClient.CreateWaitingReservationNote(ctx, &guestProto.CreateWaitingReservationNoteRequest{
			ReservationWaitlistId: reservationWaitlistId,
			Description:           input.Description,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, buildWaitingReservationNoteResponse(proto.GetResult()))
	}
}
