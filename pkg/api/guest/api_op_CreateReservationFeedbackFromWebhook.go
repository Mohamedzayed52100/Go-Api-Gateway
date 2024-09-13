package guest

import (
	"context"
	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"net/http"
)

func (s *SGuest) CreateReservationFeedbackFromWebhook() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.CreateReservationFeedbackWebhookRequest
		common.ResolveRequestBody(c, &input)

		proto, err := s.ReservationFeedbackWebhookClient.CreateReservationFeedbackFromWebhook(ctx, &guestProto.CreateReservationFeedbackFromWebhookRequest{
			ReservationId: common.ConvertStringToInt(input.ReservationId),
			Rate:          common.ConvertStringToInt(input.Rate),
			Feedback:      input.Feedback,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, &openapi.CreateReservationFeedbackWebhookResponse{
			Code:    proto.GetCode(),
			Message: proto.GetMessage(),
		})
	}
}
