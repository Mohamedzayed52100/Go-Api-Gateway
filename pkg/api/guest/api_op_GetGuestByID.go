package guest

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SGuest) GetGuestByID() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		id := c.Param("guestId")

		proto, err := s.GuestClient.GetGuestByID(ctx, &guestProto.GetGuestByIDRequest{
			Id: common.ConvertStringToInt(id),
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildGuestAPIResponse(proto.GetResult()))
	}
}
