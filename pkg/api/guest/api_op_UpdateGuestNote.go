package guest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) UpdateGuestNote() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.GuestNoteParams
		common.ResolveRequestBody(c, &input)
		noteId := common.ConvertStringToInt(c.Param("noteId"))

		proto, err := s.GuestClient.UpdateGuestNote(ctx, &guestProto.UpdateGuestNoteRequest{
			Params: &guestProto.GuestNoteParams{
				Id:          noteId,
				Description: input.Description,
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, buildGuestNoteAPIResponse(proto.GetResult()))
	}
}
