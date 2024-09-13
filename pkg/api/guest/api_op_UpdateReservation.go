package guest

import (
	"context"
	"net/http"
	"strconv"

	"github.com/goplaceapp/goplace-common/pkg/logger"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SGuest) UpdateReservation() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		id := c.Param("reservationId")
		var input openapi.UpdateReservationRequest
		common.ResolveRequestBody(c, &input)
		emptyTags := (input.Tags == nil)
		pinCode := c.GetHeader("SU-Pin-Code")

		var tags []*guestProto.TagParams
		if input.Tags == nil {
			emptyTags = true
		} else {
			for _, tag := range input.Tags {
				tags = append(tags, &guestProto.TagParams{
					Id:         tag.Id,
					CategoryId: tag.CategoryId,
				})
			}
		}

		var specialOccasionId int32
		var deleteSpecialOccasion bool
		if input.SpecialOccasionId != nil {
			specialOccasionId = *input.SpecialOccasionId
			if *input.SpecialOccasionId == 0 {
				deleteSpecialOccasion = true
			}
		}

		proto, err := s.ReservationClient.UpdateReservation(ctx, &guestProto.UpdateReservationRequest{
			Params: &guestProto.ReservationParams{
				Id:                    common.ConvertStringToInt(id),
				GuestId:               input.GuestId,
				BranchId:              input.BranchId,
				ShiftId:               input.ShiftId,
				StatusId:              input.StatusId,
				SeatingAreaId:         input.SeatingAreaId,
				GuestsNumber:          input.GuestsNumber,
				SeatedGuests:          input.SeatedGuests,
				Date:                  input.Date,
				Time:                  input.Time,
				ReservedVia:           input.ReservedVia,
				SpecialOccasionId:     specialOccasionId,
				DeleteSpecialOccasion: deleteSpecialOccasion,
				Tags:                  tags,
				EmptyTags:             emptyTags,
				PinCode:               pinCode,
			},
		})

		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		res := buildReservationAPIResponse(proto.GetResult())

		if ctx.Value(meta.TenantIDContextKey.String()) != nil {
			channelName := "reservation-" + ctx.Value(meta.TenantIDContextKey.String()).(string) + "-" + strconv.Itoa(int(res.Branch.Id)) + "-" + res.Date + "-" + strconv.Itoa(int(res.Shift.Id))
			eventName := "update-reservation"

			err := PusherClient.Trigger(channelName, eventName, res)
			if err != nil {
				logger.Default().Error("Failed to trigger Pusher event: ", err)
			}
		}

		c.JSON(http.StatusOK, res)
	}
}
