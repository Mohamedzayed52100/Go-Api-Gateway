package guest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) CreateReservationFeedback() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.CreateReservationFeedbackRequest
		common.ResolveRequestBody(c, &input)
		var (
			proto *guestProto.CreateReservationFeedbackResponse
			err   error
		)

		if input.SectionIds == nil {
			proto, err = s.ReservationFeedbackClient.CreateReservationFeedback(ctx, &guestProto.CreateReservationFeedbackRequest{
				Params: &guestProto.ReservationFeedbackParams{
					ReservationId: input.ReservationId,
					Rate:          input.Rate,
					Description:   input.Description,
				},
			})
		} else {
			proto, err = s.ReservationFeedbackClient.CreateReservationFeedback(ctx, &guestProto.CreateReservationFeedbackRequest{
				Params: &guestProto.ReservationFeedbackParams{
					ReservationId: input.ReservationId,
					SectionIds:    *input.SectionIds,
					Rate:          input.Rate,
					Description:   input.Description,
				},
			})
		}
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, buildReservationFeedbackAPIResponse(proto.GetResult()))
	}
}
