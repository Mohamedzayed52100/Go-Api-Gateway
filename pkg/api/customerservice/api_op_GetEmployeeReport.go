package customerservice

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	cssProto "github.com/goplaceapp/goplace-customer-service/api/v1"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SCustomerService) GetEmployeeReport() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		employeeId := common.ConvertStringToInt(c.Param("employeeId"))
		perPage := c.Query("perPage")
		currentPage := c.Query("currentPage")
		contentType := common.ParseQueryCriteriaToStringArray(c.Query("type"))
		inquiryStatus := common.ParseQueryCriteriaToIntArray(c.Query("inquiryStatus"))
		reservationStatus := common.ParseQueryCriteriaToIntArray(c.Query("reservationStatus"))
		fromDate := c.Query("fromDate")
		toDate := c.Query("toDate")
		query := c.Query("query")
		if len(contentType) == 0 {
			contentType = []string{"reservation", "inquiry"}
		}

		proto, err := s.CSSReportClient.GetEmployeeReport(ctx, &cssProto.GetEmployeeReportRequest{
			Params: &cssProto.CPaginationParams{
				PerPage:     common.ConvertStringToInt(perPage),
				CurrentPage: common.ConvertStringToInt(currentPage),
			},
			EmployeeId:        employeeId,
			Type:              contentType,
			InquiryStatus:     inquiryStatus,
			ReservationStatus: reservationStatus,
			FromDate:          fromDate,
			ToDate:            toDate,
			Query:             query,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, &openapi.GetEmployeeReportResponse{
			Pagination: buildPaginationResponse(proto.GetPagination()),
			Result:     buildAllEmployeeOperationsResponse(proto.GetResult()),
			Employee:   buildEmployeeResponse(proto.GetEmployee()),
		})
	}
}
