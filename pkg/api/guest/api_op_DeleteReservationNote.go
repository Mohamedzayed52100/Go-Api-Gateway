package guest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) DeleteReservationNote() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		id := common.ConvertStringToInt(c.Param("noteId"))
		reservationId := common.ConvertStringToInt(c.Param("reservationId"))
		proto, err := s.ReservationClient.DeleteReservationNote(ctx, &guestProto.DeleteReservationNoteRequest{
			Id:            id,
			ReservationId: reservationId,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, openapi.DeleteReservationNoteResponse{
			Code:    proto.GetCode(),
			Message: proto.GetMessage(),
		})
	}
}
