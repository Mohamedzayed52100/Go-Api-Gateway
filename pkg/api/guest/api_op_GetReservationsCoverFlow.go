package guest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) GetReservationsCoverFlow() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		shiftId := common.ConvertStringToInt(c.Query("shiftId"))
		seatingArea := common.ParseQueryCriteriaToIntArray(c.Query("seatingArea"))
		date := c.Query("date")

		proto, err := s.ReservationClient.GetReservationsCoverFlow(ctx, &guestProto.GetReservationsCoverFlowRequest{
			ShiftId:     shiftId,
			Date:        date,
			SeatingArea: seatingArea,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildCoverFlowResponse(proto.GetResult()))
	}
}
