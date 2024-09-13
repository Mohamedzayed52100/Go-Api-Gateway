package guest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) UpdateReservationFeedback() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.UpdateReservationFeedbackRequest
		common.ResolveRequestBody(c, &input)
		reservationId := common.ConvertStringToInt(c.Param("reservationId"))
		feedbackId := common.ConvertStringToInt(c.Param("feedbackId"))
		var (
			proto *guestProto.UpdateReservationFeedbackResponse
			err   error
		)

		if input.SectionIds == nil {
			proto, err = s.ReservationFeedbackClient.UpdateReservationFeedback(ctx, &guestProto.UpdateReservationFeedbackRequest{
				Params: &guestProto.ReservationFeedbackParams{
					Id:            feedbackId,
					Status:        input.Status,
					ReservationId: reservationId,
					EmptySections: true,
					Rate:          input.Rate,
					Description:   input.Description,
				},
			})
		} else {
			proto, err = s.ReservationFeedbackClient.UpdateReservationFeedback(ctx, &guestProto.UpdateReservationFeedbackRequest{
				Params: &guestProto.ReservationFeedbackParams{
					Id:            feedbackId,
					Status:        input.Status,
					ReservationId: reservationId,
					SectionIds:    *input.SectionIds,
					EmptySections: false,
					Rate:          input.Rate,
					Description:   input.Description,
				},
			})
		}
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildReservationFeedbackAPIResponse(proto.GetResult()))
	}
}
