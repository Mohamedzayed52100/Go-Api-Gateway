package settings

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
)

func (s *SSettings) GetAllTables() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		branch := common.ParseQueryCriteriaToIntArray(c.Query("branch"))
		proto, err := s.TableClient.GetAllTables(ctx, &settingsProto.GetAllTablesRequest{
			Branch: branch,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusAccepted, buildAllTablesResponse(proto.GetResult()))
	}
}
