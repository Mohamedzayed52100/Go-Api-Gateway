package guest

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SGuest) GetPaymentByID() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		proto, err := s.PaymentClient.GetPaymentByID(ctx, &guestProto.GetPaymentByIDRequest{
			Id:            common.ConvertStringToInt(c.Param("paymentId")),
			ReservationId: common.ConvertStringToInt(c.Param("reservationId")),
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(200, buildPaymentResponse(proto.GetResult()))
	}
}
