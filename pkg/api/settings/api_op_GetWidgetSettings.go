package settings

import (
	"context"
	"net/http"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SSettings) GetWidgetSettings() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		proto, err := s.WidgetSettingsClient.GetWidgetSettings(ctx, &emptypb.Empty{})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildWidgetSettingsResponse(proto.GetResult()))
	}
}
