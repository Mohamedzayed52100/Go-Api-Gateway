package settings

import (
	"context"
	"net/http"

	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SSettings) GetRestaurantItemByID() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		proto, err := s.RestaurantItemClient.GetRestaurantItemByID(ctx, &settingsProto.GetRestaurantItemByIDRequest{
			Id: common.ConvertStringToInt(c.Param("itemId")),
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildRestaurantItemResponse(proto.GetResult()))
	}
}
