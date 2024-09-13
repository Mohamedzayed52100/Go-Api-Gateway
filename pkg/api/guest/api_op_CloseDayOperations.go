package guest

import (
	"context"
	"net/http"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SGuest) CloseDayOperations() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		date := c.Query("date")
		pinCode := c.GetHeader("SU-Pin-Code")

		proto, err := s.DayOperationsClient.CloseDayOperations(ctx, &guestProto.CloseDayOperationsRequest{
			Date:    date,
			PinCode: pinCode,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, proto)
	}
}
