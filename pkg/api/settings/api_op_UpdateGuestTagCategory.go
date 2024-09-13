package settings

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
)

func (s *SSettings) UpdateGuestTagCategory() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		id := c.Param("categoryId")
		var input openapi.UpdateGuestTagCategoryRequest
		common.ResolveRequestBody(c, &input)

		proto, err := s.GuestTagClient.UpdateGuestTagCategory(ctx, &settingsProto.UpdateGuestTagCategoryRequest{
			Params: &settingsProto.TagCategoryParams{
				Id:             common.ConvertStringToInt(id),
				Name:           input.Name,
				Color:          input.Color,
				Classification: input.Classification,
				OrderIndex:     input.OrderIndex,
				IsDisabled:     input.IsDisabled,
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, buildTagCategoryResponse(proto.GetResult()))
	}
}
