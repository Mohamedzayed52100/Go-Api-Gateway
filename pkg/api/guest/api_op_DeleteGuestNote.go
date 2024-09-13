package guest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) DeleteGuestNote() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		id := common.ConvertStringToInt(c.Param("noteId"))
		guestId := common.ConvertStringToInt(c.Param("guestId"))
		proto, err := s.GuestClient.DeleteGuestNote(ctx, &guestProto.DeleteGuestNoteRequest{
			Id:      id,
			GuestId: guestId,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, openapi.DeleteGuestNoteResponse{
			Code:    proto.GetCode(),
			Message: proto.GetMessage(),
		})
	}
}
