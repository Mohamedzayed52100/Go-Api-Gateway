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

func (s *SGuest) CreateReservation() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.CreateReservationRequest
		common.ResolveRequestBody(c, &input)
		pinCode := c.GetHeader("SU-Pin-Code")

		var tags []*guestProto.TagParams
		for _, tag := range input.Tags {
			tags = append(tags, &guestProto.TagParams{
				Id:         tag.Id,
				CategoryId: tag.CategoryId,
			})
		}
		var sepcialOccasionId int32
		if input.SpecialOccasionId != nil {
			sepcialOccasionId = *input.SpecialOccasionId
		}

		proto, err := s.ReservationClient.CreateReservation(ctx, &guestProto.CreateReservationRequest{
			Params: &guestProto.ReservationParams{
				GuestId:           input.GuestId,
				SeatingAreaId:     input.SeatingAreaId,
				BranchId:          input.BranchId,
				ShiftId:           input.ShiftId,
				GuestsNumber:      input.GuestsNumber,
				Date:              input.Date,
				Time:              input.Time,
				CreationDuration:  input.CreationDuration,
				ReservedVia:       input.ReservedVia,
				SpecialOccasionId: sepcialOccasionId,
				Tags:              tags,
				PinCode:           pinCode,
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
	}
}
