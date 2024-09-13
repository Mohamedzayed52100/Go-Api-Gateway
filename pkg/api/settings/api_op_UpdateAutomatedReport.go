package settings

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
)

func (s *SSettings) UpdateAutomatedReport() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.UpdateAutomatedReportRequest
		common.ResolveRequestBody(c, &input)
		id := common.ConvertStringToInt(c.Param("reportId"))

		userIds := []int32{}
		if input.UserIds != nil {
			userIds = *input.UserIds
		}

		proto, err := s.AutomatedReportClient.UpdateAutomatedReport(ctx, &settingsProto.UpdateAutomatedReportRequest{
			Params: &settingsProto.AutomatedReportParams{
				Id:                  id,
				Name:                input.Name,
				TypeIds:             input.Reports,
				RepeatInterval:      input.RepeatInterval,
				StartTime:           input.StartTime,
				EndTime:             input.EndTime,
				RepeatIntervalHours: input.RepeatIntervalHours,
				UserIds:             userIds,
				EmptyUsers:          input.UserIds != nil && len(*input.UserIds) == 0,
				SendVia:             strings.Join(ListLowerCase(input.SendVia), ","),
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildAutomatedReportResponse(proto.GetResult()))
	}
}
