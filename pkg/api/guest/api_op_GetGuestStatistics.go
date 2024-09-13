package guest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) GetGuestStatistics() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		id := common.ConvertStringToInt(c.Param("guestId"))

		proto, err := s.GuestClient.GetGuestStatistics(ctx, &guestProto.GetGuestStatisticsRequest{
			GuestId: id,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildGuestStatisticsResponse(proto.GetResult()))
	}
}
