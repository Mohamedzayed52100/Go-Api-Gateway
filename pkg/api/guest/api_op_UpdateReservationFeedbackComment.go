package guest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) UpdateReservationFeedbackComment() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.UpdateReservationFeedbackCommentRequest
		common.ResolveRequestBody(c, &input)
		commentId := common.ConvertStringToInt(c.Param("commentId"))
		inquiryId := common.ConvertStringToInt(c.Param("inquiryId"))

		proto, err := s.ReservationFeedbackCommentClient.UpdateReservationFeedbackComment(ctx, &guestProto.UpdateReservationFeedbackCommentRequest{
			Params: &guestProto.ReservationFeedbackCommentParams{
				Id:                    commentId,
				ReservationFeedbackId: inquiryId,
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
