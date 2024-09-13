package settings

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *SSettings) GetAllCountries() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		proto, err := s.SettingsClient.GetAllCountries(ctx, &emptypb.Empty{})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildAllCountriesResponse(proto.GetResult()))
	}
}
