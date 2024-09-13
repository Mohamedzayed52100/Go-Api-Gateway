package settings

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
)

func (s *SSettings) CreateReservationTagCategory() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.CreateReservationTagCategoryRequest
		common.ResolveRequestBody(c, &input)

		proto, err := s.ReservationTagClient.CreateReservationTagCategory(ctx, &settingsProto.CreateReservationTagCategoryRequest{
			Params: &settingsProto.TagCategoryParams{
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
