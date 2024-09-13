package guest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) GetWidgetAllSpecialOccasions() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		proto, err := s.ReservationSpecialOccasionClient.GetWidgetAllSpecialOccasions(ctx, &guestProto.GetWidgetAllSpecialOccasionsRequest{
			BranchId: common.ConvertStringToInt(c.Query("branchId")),
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildAllSpecialOccasionsResponse(proto.GetResult()))
	}
}
