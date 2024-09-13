package guest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) AddGuestNote() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		guestId := c.Param("guestId")
		var input openapi.GuestNoteParams
		common.ResolveRequestBody(c, &input)
		proto, err := s.GuestClient.AddGuestNote(ctx, &guestProto.AddGuestNoteRequest{
			Params: &guestProto.GuestNoteParams{
				GuestId:     common.ConvertStringToInt(guestId),
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
