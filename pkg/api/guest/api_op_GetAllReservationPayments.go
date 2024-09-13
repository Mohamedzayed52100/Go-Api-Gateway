package guest

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SGuest) GetAllReservationPayments() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		proto, err := s.PaymentClient.GetAllReservationPayments(ctx, &guestProto.GetAllReservationPaymentsRequest{
			Id: common.ConvertStringToInt(c.Param("reservationId")),
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(200, buildAllPaymentsResponse(proto.GetResult()))
	}
}
