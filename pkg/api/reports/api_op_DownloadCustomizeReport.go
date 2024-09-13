package reports

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	reports "github.com/goplaceapp/shared-protobufs/reports/go_out"
)

func (s *SReports) DownloadCustomizeReport() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.DownloadCustomizeReportRequest
		common.ResolveRequestBody(c, &input)

		proto, err := s.CustomizeReportClient.DownloadCustomizeReport(ctx, &reports.DownloadReportRequest{
			Params: &reports.CustomizeReportParams{
				FromDate:            input.Filters.FromDate,
				ToDate:              input.Filters.ToDate,
				FromTime:            input.Filters.FromTime,
				ToTime:              input.Filters.ToTime,
				MinSpent:            input.Filters.MinSpent,
				MaxSpent:            input.Filters.MaxSpent,
				BranchIds:           input.Filters.BranchIds,
				ReservationTagIds:   input.Filters.ReservationTagIds,
				ShiftIds:            input.Filters.ShiftIds,
				SeatingAreaIds:      input.Filters.SeatingAreaIds,
				BookedVia:           input.Filters.BookedVia,
				ReservationStatuses: input.Filters.ReservationStatuses,
			},
			Columns: input.Columns,
			Email:   input.Email,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, &openapi.DownloadCustomizeReportResponse{
			Code:    proto.GetCode(),
			Message: proto.GetMessage(),
		})
	}
}
