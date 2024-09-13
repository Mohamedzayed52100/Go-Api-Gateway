package guest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) UpdateWaitingReservationNote() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		noteId := common.ConvertStringToInt(c.Param("noteId"))
		var input openapi.UpdateWaitingReservationNoteRequest
		common.ResolveRequestBody(c, &input)

		proto, err := s.ReservationWaitlistClient.UpdateWaitingReservationNote(ctx, &guestProto.UpdateWaitingReservationNoteRequest{
			Id:          noteId,
			Description: input.Description,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildWaitingReservationNoteResponse(proto.GetResult()))
	}
}
