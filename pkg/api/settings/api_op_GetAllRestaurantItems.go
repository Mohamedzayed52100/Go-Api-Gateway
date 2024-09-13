package settings

import (
	"context"
	"net/http"
	"strings"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	"github.com/goplaceapp/goplace-gateway/pkg/api/convert"
)

func (s *SSettings) GetAllRestaurantItems() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		sortingCriteria := common.ConvertSortingCriteria(c.Query("sort"))
		perPage := common.ConvertStringToInt(c.Query("perPage"))
		currentPage := common.ConvertStringToInt(c.Query("currentPage"))

		proto, err := s.RestaurantItemClient.GetAllRestaurantItems(ctx, &settingsProto.GetAllRestaurantItemsRequest{
			Query: c.Query("query"),
			Pagination: &settingsProto.PaginationParams{
				PerPage:     perPage,
				CurrentPage: currentPage,
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		for i, j := 0, len(sortingCriteria)-1; i < j; i, j = i+1, j-1 {
			sortingCriteria[i], sortingCriteria[j] = sortingCriteria[j], sortingCriteria[i]
		}

		for _, c := range sortingCriteria {
			fields := strings.Split(c, ":")
			switch fields[0] {
			case "code":
				SortItemsByCode(proto.GetResult(), fields[1])
			case "name":
				SortItemsByName(proto.GetResult(), fields[1])
			}
		}

		c.JSON(http.StatusOK, openapi.GetAllRestaurantItemsResponse{
			Pagination: convert.BuildPaginationResponse(&guestProto.Pagination{
				Total:       proto.GetPagination().GetTotal(),
				PerPage:     proto.GetPagination().GetPerPage(),
				CurrentPage: proto.GetPagination().GetCurrentPage(),
				LastPage:    proto.GetPagination().GetLastPage(),
				From:        proto.GetPagination().GetFrom(),
				To:          proto.GetPagination().GetTo(),
			}),
			Result: buildAllRestaurantItemsResponse(proto.GetResult()),
		})
	}
}
