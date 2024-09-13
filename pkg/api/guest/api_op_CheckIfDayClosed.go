package guest

import (
	"context"
	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"net/http"
)

func (s *SGuest) CheckIfDayClosed() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		date := c.Query("date")

		proto, err := s.DayOperationsClient.CheckIfDayClosed(ctx, &guestProto.CheckIfDayClosedRequest{
			Date: date,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		if !proto.GetClosed() {
			proto.Closed = false
		}

		c.JSON(http.StatusOK, &openapi.CheckIfDayClosedResponse{Closed: proto.GetClosed()})
	}
}
