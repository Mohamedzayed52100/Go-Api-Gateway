package reports

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	reports "github.com/goplaceapp/shared-protobufs/reports/go_out"
)

func (s *SReports) GetFinancialReport() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		fromDate := c.Query("fromDate")
		toDate := c.Query("toDate")
		profileTagIds := c.Query("profileTagIds")
		reservationTagIds := c.Query("reservationTagIds")
		bookedVia := c.Query("bookedVia")
		period := c.Query("period")
		branchIds := c.Query("branchIds")

		proto, err := s.FinancialReportClient.GetFinancialReport(ctx, &reports.ReportRequestParam{
			FromDate:          fromDate,
			ToDate:            toDate,
			ProfileTagIds:     common.ParseQueryCriteriaToIntArray(profileTagIds),
			ReservationTagIds: common.ParseQueryCriteriaToIntArray(reservationTagIds),
			BookedVia:         common.ParseQueryCriteriaToStringArray(bookedVia),
			Period:            period,
			BranchIds:         common.ParseQueryCriteriaToIntArray(branchIds),
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, &openapi.GetFinancialReportResponse{
			TotalSales:            proto.GetTotalSales(),
			AverageCheck:          proto.GetAverageCheck(),
			AverageCheckPerPerson: float32(proto.GetAverageCheckPerPerson()),
		})
	}
}
