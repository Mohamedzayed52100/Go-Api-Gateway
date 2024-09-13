package settings

import (
	"context"
	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
	"net/http"
)

func (s *SSettings) DeleteSeatingArea() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		seatingAreaId := c.Param("seatingAreaId")

		proto, err := s.SeatingAreaClient.DeleteSeatingArea(ctx, &settingsProto.DeleteSeatingAreaRequest{
			Id: common.ConvertStringToInt(seatingAreaId),
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, &openapi.DeleteSeatingAreaResponse{Result: proto.GetResult()})
	}
}
