package guest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) CreateReservationFeedbackComment() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.CreateReservationFeedbackCommentRequest
		common.ResolveRequestBody(c, &input)
		reservationFeedbackId := common.ConvertStringToInt(c.Param("feedbackId"))

		proto, err := s.ReservationFeedbackCommentClient.CreateReservationFeedbackComment(ctx, &guestProto.CreateReservationFeedbackCommentRequest{
			Params: &guestProto.ReservationFeedbackCommentParams{
				ReservationFeedbackId: reservationFeedbackId,
				Comment:               input.Comment,
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, buildReservationFeedbackCommentResponse(proto.GetResult()))
	}
}
