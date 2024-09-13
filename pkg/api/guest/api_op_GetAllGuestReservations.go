package guest

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	"github.com/goplaceapp/goplace-gateway/pkg/api/convert"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) GetAllGuestReservations() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		guestID := c.Param("guestId")
		perPage := c.Query("perPage")
		currentPage := c.Query("currentPage")
		sortingCriteria := common.ConvertSortingCriteria(c.Query("sort"))

		proto, err := s.GuestClient.GetAllGuestReservations(ctx, &guestProto.GetAllGuestReservationsRequest{
			GuestId: common.ConvertStringToInt(guestID),
			Params: &guestProto.PaginationParams{
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
			case "guestsNumber":
				SortReservationsByGuestsNumber(proto.GetResult(), fields[1])
			case "time":
				SortReservationsByTime(proto.GetResult(), fields[1])
			}
		}

		c.JSON(http.StatusOK, &openapi.GetAllReservationsResponse{
			Pagination: convert.BuildPaginationResponse(proto.GetPagination()),
			Result:     buildReservationsAPIResponse(proto.GetResult()),
		})
	}
}
