package guest

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goplaceapp/goplace-common/pkg/logger"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) DeleteReservation() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		id := common.ConvertStringToInt(c.Param("reservationId"))

		proto, err := s.ReservationClient.DeleteReservation(ctx, &guestProto.DeleteReservationRequest{
			Id: id,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		if ctx.Value(meta.TenantIDContextKey.String()) != nil {
			channelName := "reservation-" + ctx.Value(meta.TenantIDContextKey.String()).(string) + "-" + strconv.Itoa(int(proto.Result.Branch.Id)) + "-" + proto.Result.Date.AsTime().Format(time.DateOnly) + "-" + strconv.Itoa(int(proto.Result.Shift.Id))
			eventName := "delete-reservation"

			err := PusherClient.Trigger(channelName, eventName, map[string]interface{}{
				"id": id,
			})
			if err != nil {
				logger.Default().Error("Failed to trigger pusher event", err)
			}
		}

		c.JSON(http.StatusOK, &openapi.DeleteReservationResponse{
			Code:    http.StatusOK,
			Message: "Reservation deleted successfully",
		})
	}
}
