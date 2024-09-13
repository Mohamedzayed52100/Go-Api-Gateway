package guest

import (
	"context"
	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"net/http"
)

func (s *SGuest) GetAvailableTimes() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		branchId := c.Param("branchId")
		partySize := c.Query("partySize")
		shiftId := c.Query("shiftId")
		seatingAreaId := c.Query("seatingAreaId")
		date := c.Query("date")

		proto, err := s.ReservationClient.GetAvailableTimes(ctx, &guestProto.GetAvailableTimesRequest{
			BranchId:      common.ConvertStringToInt(branchId),
			Date:          date,
			PartySize:     common.ConvertStringToInt(partySize),
			ShiftId:       common.ConvertStringToInt(shiftId),
			SeatingAreaId: common.ConvertStringToInt(seatingAreaId),
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildAvailableTimesResponse(proto.GetAvailableTimes()))
	}
}
