package settings

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
)

func (s *SSettings) GetAllSeatingAreas() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		branchId := common.ConvertStringToInt(c.Query("branchId"))
		proto, err := s.SeatingAreaClient.GetAllSeatingAreas(ctx, &settingsProto.GetAllSeatingAreasRequest{
			BranchId: branchId,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusAccepted, buildAllSeatingAreasResponse(proto.GetResult()))
	}
}
