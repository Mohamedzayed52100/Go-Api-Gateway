package customerservice

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	cssProto "github.com/goplaceapp/goplace-customer-service/api/v1"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SCustomerService) GetAllEmployeesReports() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		perPage := c.Query("perPage")
		currentPage := c.Query("currentPage")
		branchIds := common.ParseQueryCriteriaToIntArray(c.Query("branch"))
		fromDate := c.Query("fromDate")
		toDate := c.Query("toDate")
		query := c.Query("query")

		proto, err := s.CSSReportClient.GetAllEmployeesReport(ctx, &cssProto.GetAllEmployeesReportRequest{
			Params: &cssProto.CPaginationParams{
				PerPage:     common.ConvertStringToInt(perPage),
				CurrentPage: common.ConvertStringToInt(currentPage),
			},
			Branch:   branchIds,
			FromDate: fromDate,
			ToDate:   toDate,
			Query:    query,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		res := &openapi.GetAllEmployeesReportsResponse{
			Pagination: buildPaginationResponse(proto.GetPagination()),
			Result:     buildAllEmployeesReportsResponse(proto.GetResult()),
			Statistics: buildStatisticsResponse(proto.GetStatistics()),
		}

		c.JSON(http.StatusOK, res)
	}
}
