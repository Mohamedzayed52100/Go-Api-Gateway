package customerservice

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	cssProto "github.com/goplaceapp/goplace-customer-service/api/v1"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SCustomerService) GetAllInquiries() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		perPage := c.Query("perPage")
		currentPage := c.Query("currentPage")
		query := c.Query("query")
		branchIds := common.ParseQueryCriteriaToIntArray(c.Query("branch"))
		statusIds := common.ParseQueryCriteriaToIntArray(c.Query("status"))
		fromDate := c.Query("fromDate")
		toDate := c.Query("toDate")
		reservationId := c.Query("reservationId")

		proto, err := s.CSSInquiryClient.GetAllInquiries(ctx, &cssProto.GetAllInquiriesRequest{
			Params: &cssProto.CPaginationParams{
				PerPage:     common.ConvertStringToInt(perPage),
				CurrentPage: common.ConvertStringToInt(currentPage),
			},
			Query:         query,
			Branch:        branchIds,
			Status:        statusIds,
			FromDate:      fromDate,
			ToDate:        toDate,
			ReservationId: common.ConvertStringToInt(reservationId),
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, &openapi.GetAllInquiriesResponse{
			Pagination:   buildPaginationResponse(proto.GetPagination()),
			Result:       buildAllInquiriesResponse(proto.GetResult()),
			TotalPending: int32(proto.GetTotalPending()),
			TotalSolved:  int32(proto.GetTotalSolved()),
		})
	}
}
