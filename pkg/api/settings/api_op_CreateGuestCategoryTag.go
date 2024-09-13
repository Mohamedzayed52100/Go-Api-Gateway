package settings

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
)

func (s *SSettings) CreateGuestCategoryTag() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.CreateGuestCategoryTagRequest
		common.ResolveRequestBody(c, &input)

		proto, err := s.GuestTagClient.CreateGuestTag(ctx, &settingsProto.CreateGuestTagRequest{
			Params: &settingsProto.TagParams{
				Name:       input.Name,
				CategoryId: input.CategoryId,
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, buildTagResponse(proto.GetResult()))
	}
}
