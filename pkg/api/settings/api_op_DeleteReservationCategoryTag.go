package settings

import (
	"context"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
)

func (s *SSettings) DeleteReservationCategoryTag() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		id := common.ConvertStringToInt(c.Param("tagId"))

		proto, err := s.ReservationTagClient.DeleteReservationTag(ctx, &settingsProto.DeleteReservationTagRequest{
			Id: id,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(200, openapi.DeleteTagResponse{
			Code:    proto.GetCode(),
			Message: proto.GetMessage(),
		})
	}
}
