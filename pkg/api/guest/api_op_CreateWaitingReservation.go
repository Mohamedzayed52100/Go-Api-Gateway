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

func (s *SGuest) CreateWaitingReservation() openapi.ContextHandler {
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

		switch input.Type {
		case "waiting":
			proto, err := s.ReservationWaitlistClient.CreateWaitingReservation(ctx, &guestProto.CreateWaitingReservationRequest{
				Params: &guestProto.ReservationWaitlistParams{
					GuestId:       input.GuestId,
					SeatingAreaId: input.SeatingAreaId,
					ShiftId:       input.ShiftId,
					GuestsNumber:  input.GuestsNumber,
					WaitingTime:   input.WaitingTime,
					Tags:          tags,
					NoteId:        input.NoteId,
					Type:          input.Type,
					Date:          input.Date,
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
				eventName := "new-waiting-reservation"

				err = PusherClient.Trigger(channelName, eventName, res)
				if err != nil {
					logger.Default().Error("Failed to trigger Pusher event: ", err)
				}
			}

			c.JSON(http.StatusCreated, res)
			return
		case "direct-in":
			proto, err := s.ReservationWaitlistClient.CreateDirectInReservation(ctx, &guestProto.CreateWaitingReservationRequest{
				Params: &guestProto.ReservationWaitlistParams{
					GuestId:       input.GuestId,
					SeatingAreaId: input.SeatingAreaId,
					ShiftId:       input.ShiftId,
					GuestsNumber:  input.GuestsNumber,
					Tags:          tags,
					Date:          input.Date,
				},
			})
			if err != nil {
				common.HandleGrpcError(c, err)
				return
			}

			res := buildReservationAPIResponse(proto.GetResult())

			if ctx.Value(meta.TenantIDContextKey.String()) != nil {
				channelName := "reservation-" + ctx.Value(meta.TenantIDContextKey.String()).(string) + "-" + strconv.Itoa(int(res.Branch.Id)) + "-" + res.Date + "-" + strconv.Itoa(int(res.Shift.Id))
				eventName := "new-reservation"

				err = PusherClient.Trigger(channelName, eventName, res)
				if err != nil {
					logger.Default().Error("Failed to trigger Pusher event: ", err)
				}
			}

			c.JSON(http.StatusCreated, res)
			return
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation type"})
			return
		}
	}
}
