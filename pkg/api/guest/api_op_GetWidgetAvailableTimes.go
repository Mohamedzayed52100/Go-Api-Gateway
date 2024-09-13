package guest

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) GetWidgetAvailableTimes() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		branchId := c.Param("branchId")
		partySize := c.Query("partySize")
		seatingAreaId := c.Query("seatingAreaId")
		fromDate := c.Query("fromDate")
		toDate := c.Query("toDate")

		switch {
		case branchId == "":
			c.JSON(http.StatusBadRequest, &openapi.GeneralError{
				Code:    http.StatusBadRequest,
				Message: "branchId param is required",
			})
			return
		case partySize == "":
			c.JSON(http.StatusBadRequest, &openapi.GeneralError{
				Code:    http.StatusBadRequest,
				Message: "partySize param is required",
			})
			return
		case seatingAreaId == "":
			c.JSON(http.StatusBadRequest, &openapi.GeneralError{
				Code:    http.StatusBadRequest,
				Message: "seatingAreaId param is required",
			})
			return
		case fromDate == "":
			fromDate = time.Now().Format(time.DateOnly)
		case toDate == "":
			toDate = time.Now().Format(time.DateOnly)
		}

		proto, err := s.ReservationWidgetClient.GetWidgetAvailableTimes(ctx, &guestProto.GetWidgetAvailableTimesRequest{
			BranchId:      common.ConvertStringToInt(branchId),
			FromDate:      fromDate,
			ToDate:        toDate,
			PartySize:     common.ConvertStringToInt(partySize),
			SeatingAreaId: common.ConvertStringToInt(seatingAreaId),
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildWidgetAvailableTimesResponse(proto.GetAvailableTimes()))
	}
}
