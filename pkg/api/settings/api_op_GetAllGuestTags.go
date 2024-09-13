package settings

import (
	"context"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SSettings) GetAllGuestTags() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		proto, err := s.GuestTagClient.GetAllGuestTags(ctx, &settingsProto.GetAllGuestTagsRequest{
			Query: c.Query("query"),
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		result := proto.GetResult()
		if len(result) == 0 {
			c.JSON(http.StatusAccepted, []openapi.Tag{})
		} else {
			c.JSON(http.StatusAccepted, buildAllTagsResponse(result))
		}
	}
}
