package reports

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	reports "github.com/goplaceapp/shared-protobufs/reports/go_out"
)

func (s *SReports) GetCustomizeReport() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var (
			fromDate            = c.Query("fromDate")
			toDate              = c.Query("toDate")
			fromTime            = c.Query("fromTime")
			toTime              = c.Query("toTime")
			bookedVia           = common.ParseQueryCriteriaToStringArray(c.Query("bookedVia"))
			branchIds           = common.ParseQueryCriteriaToIntArray(c.Query("branchIds"))
			tagIds              = common.ParseQueryCriteriaToIntArray(c.Query("tagIds"))
			shiftIds            = common.ParseQueryCriteriaToIntArray(c.Query("shiftIds"))
			seatingAreaIds      = common.ParseQueryCriteriaToIntArray(c.Query("seatingAreaIds"))
			reservationStatuses = common.ParseQueryCriteriaToStringArray(c.Query("reservationStatuses"))
			minSpent            = common.ConvertStringToInt(c.Query("minSpent"))
			maxSpent            = common.ConvertStringToInt(c.Query("maxSpent"))
			perPage             = common.ConvertStringToInt(c.Query("perPage"))
			currentPage         = common.ConvertStringToInt(c.Query("currentPage"))
		)

		proto, err := s.CustomizeReportClient.GetCustomizeReport(ctx, &reports.CustomizeReportRequest{
			Params: &reports.CustomizeReportParams{
				FromDate:            fromDate,
				ToDate:              toDate,
				FromTime:            fromTime,
				ToTime:              toTime,
				BookedVia:           bookedVia,
				BranchIds:           branchIds,
				ReservationTagIds:   tagIds,
				ShiftIds:            shiftIds,
				SeatingAreaIds:      seatingAreaIds,
				ReservationStatuses: reservationStatuses,
				MinSpent:            minSpent,
				MaxSpent:            maxSpent,
			},
			PerPage:     perPage,
			CurrentPage: currentPage,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		res := &openapi.GetCustomizeReportResponse{
			Result:     buildCustomizeReportResponse(proto.GetResult()),
			Pagination: BuildPaginationResponse(proto.Pagination),
		}

		if proto.GetResult() == nil {
			res.Result = []openapi.CustomizeReport{}
		}

		c.JSON(http.StatusOK, res)
	}
}
