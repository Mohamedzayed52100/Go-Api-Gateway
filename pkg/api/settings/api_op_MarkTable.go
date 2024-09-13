package settings

import (
	"context"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
)

func (s *SSettings) MarkTable() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.MarkTableRequest
		common.ResolveRequestBody(c, &input)
		id := common.ConvertStringToInt(c.Param("id"))

		proto, err := s.TableClient.MarkTable(ctx, &settingsProto.MarkTableRequest{
			TableId:          id,
			Status:           input.Status,
			MoveReservations: input.MoveReservations,
			Date:             input.Date,
		})

		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(200, buildTableResponse(proto.GetResult()))
	}
}
