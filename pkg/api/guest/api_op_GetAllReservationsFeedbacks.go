package guest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	"github.com/goplaceapp/goplace-gateway/pkg/api/convert"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) GetAllReservationsFeedbacks() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		perPage := common.ConvertStringToInt(c.Query("perPage"))
		currentPage := common.ConvertStringToInt(c.Query("currentPage"))
		query := c.Query("query")
		fromDate := c.Query("fromDate")
		toDate := c.Query("toDate")
		branchIds := parseQueryCriteriaToIntArray(c.Query("branch"))
		statusIds := parseQueryCriteriaToIntArray(c.Query("status"))
		rate := parseQueryCriteriaToStringArray(c.Query("rate"))

		proto, err := s.ReservationFeedbackClient.GetAllReservationsFeedbacks(ctx, &guestProto.GetAllReservationsFeedbacksRequest{
			PaginationParams: &guestProto.PaginationParams{
				PerPage:     perPage,
				CurrentPage: currentPage,
			},
			Query:     query,
			FromDate:  fromDate,
			ToDate:    toDate,
			BranchIds: branchIds,
			StatusIds: statusIds,
			Rate:      rate,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, &openapi.GetAllReservationsFeedbacksResponse{
			Pagination:    convert.BuildPaginationResponse(proto.GetPagination()),
			Result:        buildAllReservationFeedbacksResponse(proto.GetResult()),
			TotalPositive: proto.GetTotalPositive(),
			TotalNegative: proto.GetTotalNegative(),
			TotalPending:  proto.GetTotalPending(),
			TotalSolved:   proto.GetTotalSolved(),
		})
	}
}
