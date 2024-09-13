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

func (s *SGuest) AddReservationNote() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.ReservationNoteParams
		common.ResolveRequestBody(c, &input)
		id := common.ConvertStringToInt(c.Param("reservationId"))

		proto, err := s.ReservationClient.AddReservationNote(ctx, &guestProto.AddReservationNoteRequest{
			Params: &guestProto.ReservationNoteParams{
				ReservationId: id,
				Description:   input.Description,
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		res := buildReservationAPIResponse(proto.GetResult().GetReservation())

		if ctx.Value(meta.TenantIDContextKey.String()) != nil {
			channelName := "reservation-" + ctx.Value(meta.TenantIDContextKey.String()).(string) + "-" + strconv.Itoa(int(res.Branch.Id)) + "-" + res.Date + "-" + strconv.Itoa(int(res.Shift.Id))
			eventName := "update-reservation"

			err := PusherClient.Trigger(channelName, eventName, res)
			if err != nil {
				logger.Default().Error("Failed to trigger Pusher event: ", err)
			}
		}

		c.JSON(http.StatusCreated, buildReservationNoteAPIResponse(proto.GetResult()))
	}
}
