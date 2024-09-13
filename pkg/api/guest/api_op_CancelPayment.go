package guest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) CancelPayment() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		proto, err := s.PaymentClient.CancelPayment(ctx, &guestProto.CancelPaymentRequest{
			InvoiceId: common.ConvertStringToInt(c.Param("invoiceId")),
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, openapi.CancelPaymentResponse{
			Code:    proto.GetCode(),
			Message: proto.GetMessage(),
		})
	}
}
