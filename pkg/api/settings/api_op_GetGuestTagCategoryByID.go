package settings

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
)

func (s *SSettings) GetGuestTagCategoryByID() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		id := c.Param("categoryId")

		proto, err := s.GuestTagClient.GetGuestTagCategoryByID(ctx, &settingsProto.GetGuestTagCategoryByIDRequest{
			Id: common.ConvertStringToInt(id),
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusAccepted, buildTagCategoryResponse(proto.GetResult()))
	}
}
