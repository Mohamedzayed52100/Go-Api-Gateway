package settings

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
)

func (s *SSettings) GetAllReservationTags() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		proto, err := s.ReservationTagClient.GetAllReservationTags(ctx, &settingsProto.GetAllReservationTagsRequest{
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
