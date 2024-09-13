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

func (s *SGuest) UpdateWaitingReservationDetails() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.CreateWaitingReservationRequest
		common.ResolveRequestBody(c, &input)

		var tags []*guestProto.TagParams
		if input.Tags != nil {
			for _, tag := range *input.Tags {
				tags = append(tags, &guestProto.TagParams{
					Id:         tag.Id,
					CategoryId: tag.CategoryId,
				})
			}
		}

		id := common.ConvertStringToInt(c.Param("reservationWaitlistId"))

		proto, err := s.ReservationWaitlistClient.UpdateWaitingReservationDetails(ctx, &guestProto.UpdateWaitingReservationDetailsRequest{
			Params: &guestProto.ReservationWaitlistParams{
				Id:            id,
				GuestId:       input.GuestId,
				SeatingAreaId: input.SeatingAreaId,
				ShiftId:       input.ShiftId,
				GuestsNumber:  input.GuestsNumber,
				WaitingTime:   input.WaitingTime,
				Tags:          tags,
				DeleteTags:    (input.Tags != nil),
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		res := buildReservationWaitlistResponse(proto.GetResult())

		if ctx.Value(meta.TenantIDContextKey.String()) != nil {
			parseDate, err := time.Parse(time.RFC3339, proto.Result.Date)
			if err != nil {
				logger.Default().Error("Failed to parse date: ", err)
			}

			channelName := "waitlist-" + ctx.Value(meta.TenantIDContextKey.String()).(string) + "-" + strconv.Itoa(int(proto.Result.BranchId)) + "-" + parseDate.Format(time.DateOnly) + "-" + strconv.Itoa(int(proto.Result.Shift.Id))
			eventName := "update-waiting-reservation"

			err = PusherClient.Trigger(channelName, eventName, res)
			if err != nil {
				logger.Default().Error("Failed to trigger Pusher event: ", err)
			}
		}

		c.JSON(http.StatusCreated, res)
	}
}
