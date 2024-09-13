package guest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) AddReservationFeedbackSolution() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.AddReservationFeedbackSolutionRequest
		common.ResolveRequestBody(c, &input)
		feedbackId := common.ConvertStringToInt(c.Param("feedbackId"))

		proto, err := s.ReservationFeedbackCommentClient.AddReservationFeedbackSolution(ctx, &guestProto.AddReservationFeedbackSolutionRequest{
			Params: &guestProto.ReservationFeedbackSolutionParams{
				FeedbackId: feedbackId,
				Solution:   input.Solution,
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, buildReservationFeedbackSolutionResponse(proto.GetResult()))
	}
}
