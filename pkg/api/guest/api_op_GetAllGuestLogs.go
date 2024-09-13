package guest

import (
	"context"
	"net/http"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SGuest) GetAllGuestLogs() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		guestId := c.Param("guestId")

		proto, err := s.GuestLogClient.GetAllGuestLogs(ctx, &guestProto.GetAllGuestLogsRequest{
			GuestId: common.ConvertStringToInt(guestId),
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildAllGuestLogsResponse(proto.GetResult()))
	}
}
