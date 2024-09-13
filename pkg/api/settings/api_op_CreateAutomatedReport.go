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

func (s *SSettings) CreateAutomatedReport() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.CreateAutomatedReportRequest
		common.ResolveRequestBody(c, &input)

		userIds := []int32{}
		if input.UserIds != nil {
			userIds = *input.UserIds
		}

		proto, err := s.AutomatedReportClient.CreateAutomatedReport(ctx, &settingsProto.CreateAutomatedReportRequest{
			Params: &settingsProto.AutomatedReportParams{
				Name:                input.Name,
				TypeIds:             input.Reports,
				RepeatInterval:      input.RepeatInterval,
				StartTime:           input.StartTime,
				EndTime:             input.EndTime,
				RepeatIntervalHours: input.RepeatIntervalHours,
				UserIds:             userIds,
				SendVia:             strings.Join(ListLowerCase(input.SendVia), ","),
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, buildAutomatedReportResponse(proto.GetResult()))
	}
}
