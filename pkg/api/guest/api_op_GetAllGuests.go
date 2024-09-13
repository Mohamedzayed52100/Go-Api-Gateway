package guest

import (
	"context"
	"net/http"
	"strings"

	"github.com/goplaceapp/goplace-gateway/pkg/api/convert"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SGuest) GetAllGuests() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		perPage := c.Query("perPage")
		currentPage := c.Query("currentPage")
		sortingCriteria := common.ConvertSortingCriteria(c.Query("sort"))
		query := c.Query("query")

		proto, err := s.GuestClient.GetAllGuests(ctx, &guestProto.GetAllGuestsRequest{
			Query: query,
			PaginationParams: &guestProto.PaginationParams{
				PerPage:     common.ConvertStringToInt(perPage),
				CurrentPage: common.ConvertStringToInt(currentPage),
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
			case "guestName":
				SortGuestsByName(proto.GetResult(), fields[1])
			case "lastVisit":
				SortGuestsByLastVisit(proto.GetResult(), fields[1])
			}
		}

		c.JSON(http.StatusOK, &openapi.GetAllGuestsResponse{
			Pagination: convert.BuildPaginationResponse(proto.GetPagination()),
			Result:     buildGuestsAPIResponse(proto.GetResult()),
		})
	}
}
