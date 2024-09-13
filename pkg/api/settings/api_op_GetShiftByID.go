package settings

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
)

func (s *SSettings) GetShiftByID() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		id := c.Param("shiftId")

		proto, err := s.ShiftClient.GetShiftByID(ctx, &settingsProto.GetShiftByIDRequest{
			Id: common.ConvertStringToInt(id),
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusAccepted, buildShiftResponse(proto.GetResult()))
	}
}
