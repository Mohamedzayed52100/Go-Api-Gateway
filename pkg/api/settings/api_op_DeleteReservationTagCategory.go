package settings

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
)

func (s *SSettings) DeleteReservationTagCategory() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		id := common.ConvertStringToInt(c.Param("categoryId"))
		proto, err := s.ReservationTagClient.DeleteReservationTagCategory(ctx, &settingsProto.DeleteReservationTagCategoryRequest{
			Id: id,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, openapi.DeleteReservationTagCategoryResponse{
			Code:    proto.GetCode(),
			Message: proto.GetMessage(),
		})
	}
}
