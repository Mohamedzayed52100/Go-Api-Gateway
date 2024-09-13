package guest

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/goplaceapp/goplace-common/pkg/logger"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) SeatWaitingReservation() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var (
			reservationWaitlistID = c.Param("reservationWaitlistId")
			input                 openapi.UpdateWaitingReservationRequest
		)

		common.ResolveRequestBody(c, &input)

		proto, err := s.ReservationWaitlistClient.SeatWaitingReservation(ctx, &guestProto.SeatWaitingReservationRequest{
			Id:     common.ConvertStringToInt(reservationWaitlistID),
			Status: input.StatusID,
			Tables: input.Tables,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
		}

		res := buildReservationAPIResponse(proto.GetResult())

		if ctx.Value(meta.TenantIDContextKey.String()) != nil {
			// Trigger Pusher delete waitin reservation event
			channelName := "waitlist-" + ctx.Value(meta.TenantIDContextKey.String()).(string) + "-" + strconv.Itoa(int(res.Branch.Id)) + "-" + res.Date + "-" + strconv.Itoa(int(res.Shift.Id))
			eventName := "delete-waiting-reservation"
			err = PusherClient.Trigger(channelName, eventName, map[string]interface{}{
				"id": common.ConvertStringToInt(reservationWaitlistID),
			})
			if err != nil {
				logger.Default().Error("Failed to trigger Pusher event: ", err)
			}

			// Trigger Pusher new reservation event
			channelName = "reservation-" + ctx.Value(meta.TenantIDContextKey.String()).(string) + "-" + strconv.Itoa(int(res.Branch.Id)) + "-" + res.Date + "-" + strconv.Itoa(int(res.Shift.Id))
			eventName = "new-reservation"
			err = PusherClient.Trigger(channelName, eventName, res)
			if err != nil {
				logger.Default().Error("Failed to trigger Pusher event: ", err)
			}
		}

		c.JSON(http.StatusOK, res)
	}
}
