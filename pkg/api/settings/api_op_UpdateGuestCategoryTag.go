package settings

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
)

func (s *SSettings) UpdateGuestCategoryTag() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.UpdateGuestCategoryTagRequest
		common.ResolveRequestBody(c, &input)
		id := common.ConvertStringToInt(c.Param("tagId"))

		proto, err := s.GuestTagClient.UpdateGuestTag(ctx, &settingsProto.UpdateGuestTagRequest{
			Params: &settingsProto.TagParams{
				Id:         id,
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
