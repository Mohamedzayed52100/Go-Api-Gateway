package settings

import (
	"context"
	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
	"net/http"
)

func (s *SSettings) UpdateSeatingArea() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		seatingAreaId := c.Param("seatingAreaId")

		var input openapi.UpdateSeatingAreaRequest
		common.ResolveRequestBody(c, &input)

		proto, err := s.SeatingAreaClient.UpdateSeatingArea(ctx, &settingsProto.UpdateSeatingAreaRequest{
			Params: &settingsProto.SeatingAreaParams{
				Id:   common.ConvertStringToInt(seatingAreaId),
				Name: input.Name,
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildSeatingAreaResponse(proto.GetResult()))
	}
}
