package settings

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	"github.com/goplaceapp/goplace-gateway/pkg/api/convert"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
)

func (s *SSettings) GetAllShifts() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		fromDate := c.Query("fromDate")
		toDate := c.Query("toDate")
		perPage := c.Query("perPage")
		currentPage := c.Query("currentPage")
		active := c.Query("active")
		branchId := common.ConvertStringToInt(c.Query("branchId"))

		proto, err := s.ShiftClient.GetAllShifts(ctx, &settingsProto.GetAllShiftsRequest{
			FromDate: fromDate,
			ToDate:   toDate,
			Pagination: &settingsProto.PaginationParams{
				PerPage:     common.ConvertStringToInt(perPage),
				CurrentPage: common.ConvertStringToInt(currentPage),
			},
			Active:   common.ConvertStringToBool(active),
			BranchId: branchId,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusAccepted, &openapi.GetAllShiftsResponse{
			Pagination: convert.BuildPaginationResponse(&guestProto.Pagination{
				Total:       proto.GetPagination().GetTotal(),
				PerPage:     proto.GetPagination().GetPerPage(),
				CurrentPage: proto.GetPagination().GetCurrentPage(),
				LastPage:    proto.GetPagination().GetLastPage(),
				From:        proto.GetPagination().GetFrom(),
				To:          proto.GetPagination().GetTo(),
			}),
			Result: buildAllShiftsResponse(proto.GetResult()),
		})
	}
}
