package guest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) GetGuestReservationStatistics() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		id := common.ConvertStringToInt(c.Param("guestId"))
		fromDate := c.Query("fromDate")
		toDate := c.Query("toDate")
		queryType := c.Query("queryType")

		proto, err := s.GuestClient.GetGuestReservationStatistics(ctx, &guestProto.GetGuestReservationStatisticsRequest{
			GuestId:   id,
			FromDate:  fromDate,
			ToDate:    toDate,
			QueryType: queryType,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildGuestReservationStatisticsResponse(proto.GetResult()))
	}
}
