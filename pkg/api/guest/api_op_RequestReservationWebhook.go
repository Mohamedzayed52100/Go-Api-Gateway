package guest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) RequestReservationWebhook() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.RequestReservationWebhookRequest
		common.ResolveRequestBody(c, &input)

		_, err := s.ReservationClient.RequestReservationWebhook(ctx, &guestProto.RequestReservationWebhookRequest{
			PhoneNumber: input.PhoneNumber,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, openapi.RequestReservationWebhookResponse{
			Code:    http.StatusOK,
			Message: "Reservation Request has been sent successfully",
		})
	}
}
