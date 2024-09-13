package guest

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"google.golang.org/grpc/metadata"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SGuest) UpdatePaymentFromWebhook() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.UpdatePaymentFromWebhookRequest
		common.ResolveRequestBody(c, &input)

		ctx = metadata.NewOutgoingContext(ctx, metadata.New(
			map[string]string{
				"clientId": c.Param("clientId"),
			},
		))

		proto, err := s.PaymentClient.UpdatePaymentFromWebhook(ctx, &guestProto.UpdatePaymentFromWebhookRequest{
			InvoiceId: input.Id,
			Status:    input.Status,
			Card: &guestProto.CardInfo{
				FourDigits: input.Transactions[0].Card.Last_four,
				Brand:      input.Transactions[0].Card.Brand,
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(200, openapi.UpdatePaymentFromWebhookResponse{
			Code:    proto.GetCode(),
			Message: proto.GetMessage(),
		})
	}
}
