package guest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) GetAllReservationFeedbackComments() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		feedbackId := common.ConvertStringToInt(c.Param("feedbackId"))
		proto, err := s.ReservationFeedbackCommentClient.GetAllReservationFeedbackComments(ctx, &guestProto.GetAllReservationFeedbackCommentsRequest{
			ReservationFeedbackId: feedbackId,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, buildAllReservationFeedbackCommentsResponse(proto.GetResult()))
	}
}
