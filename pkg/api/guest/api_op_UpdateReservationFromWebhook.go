package guest

import (
	"context"
	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"net/http"
)

func (s *SGuest) UpdateReservationFromWebhook() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.UpdateReservationFromWebhookRequest
		common.ResolveRequestBody(c, &input)

		proto, err := s.ReservationClient.UpdateReservationFromWebhook(ctx, &guestProto.UpdateReservationFromWebhookRequest{
			ReservationId: common.ConvertStringToInt(input.ReservationId),
			Status:        input.Status,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, &openapi.UpdateReservationFromWebhookResponse{
			Code:    proto.GetCode(),
			Message: proto.GetMessage(),
		})
	}
}
