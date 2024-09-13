package reports

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	reports "github.com/goplaceapp/shared-protobufs/reports/go_out"
)

func (s *SReports) GetTopSpendersFinancialReport() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		fromDate := c.Query("fromDate")
		toDate := c.Query("toDate")
		profileTagIds := c.Query("profileTagIds")
		reservationTagIds := c.Query("reservationTagIds")
		bookedVia := c.Query("bookedVia")
		period := c.Query("period")
		perPage := common.ConvertStringToInt(c.Query("perPage"))
		currentPage := common.ConvertStringToInt(c.Query("currentPage"))
		branchIds := c.Query("branchIds")

		proto, err := s.FinancialReportClient.GetTopSpenders(ctx, &reports.TopSpenderRequestParam{
			Params: &reports.ReportRequestParam{
				FromDate:          fromDate,
				ToDate:            toDate,
				ProfileTagIds:     common.ParseQueryCriteriaToIntArray(profileTagIds),
				ReservationTagIds: common.ParseQueryCriteriaToIntArray(reservationTagIds),
				BookedVia:         common.ParseQueryCriteriaToStringArray(bookedVia),
				Period:            period,
				BranchIds:         common.ParseQueryCriteriaToIntArray(branchIds),
			},
			PerPage:     perPage,
			CurrentPage: currentPage,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, &openapi.GetFinancialTopSpendersReportResponse{
			Result:     buildTopSpendersFinancialReportResponse(proto),
			Pagination: BuildPaginationResponse(proto.GetPagination()),
		})
	}
}
