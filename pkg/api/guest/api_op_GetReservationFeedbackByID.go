package guest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) GetReservationFeedbackByID() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		reservationId := common.ConvertStringToInt(c.Param("reservationId"))
		feedbackId := common.ConvertStringToInt(c.Param("feedbackId"))

		proto, err := s.ReservationFeedbackClient.GetReservationFeedbackByID(ctx, &guestProto.GetReservationFeedbackByIDRequest{
			FeedbackId:    feedbackId,
			ReservationId: reservationId,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, buildReservationFeedbackAPIResponse(proto.GetResult()))
	}
}
