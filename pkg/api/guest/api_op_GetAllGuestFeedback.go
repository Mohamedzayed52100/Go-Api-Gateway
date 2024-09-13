package guest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) GetAllGuestFeedback() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		proto, err := s.GuestClient.GetAllGuestFeedback(ctx, &guestProto.GetAllGuestFeedbackRequest{
			GuestId: common.ConvertStringToInt(c.Param("guestId")),
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildAllReservationFeedbacksResponse(proto.GetResult()))
	}
}
