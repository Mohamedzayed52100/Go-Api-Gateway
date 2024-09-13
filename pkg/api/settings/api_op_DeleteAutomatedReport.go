package settings

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
)

func (s *SSettings) DeleteAutomatedReport() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		id := common.ConvertStringToInt(c.Param("reportId"))
		proto, err := s.AutomatedReportClient.DeleteAutomatedReport(ctx, &settingsProto.DeleteAutomatedReportRequest{
			Id: id,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, openapi.DeleteAutomatedReportResponse{
			Code:    proto.GetCode(),
			Message: proto.GetMessage(),
		})
	}
}
