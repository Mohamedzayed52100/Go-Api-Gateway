package settings

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
)

func (s *SSettings) GetCountryByName() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		countryName := c.Param("countryName")

		proto, err := s.SettingsClient.GetCountryByName(ctx, &settingsProto.GetCountryByNameRequest{
			Name: countryName,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusAccepted, buildCountryResponse(proto.GetResult()))
	}
}
