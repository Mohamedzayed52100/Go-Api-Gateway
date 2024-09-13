package guest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) CreateWidgetReservation() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.CreateWidgetReservationRequest
		common.ResolveRequestBody(c, &input)

		guests := []*guestProto.WidgetGuestParams{}
		for _, guest := range input.Guests {
			guests = append(guests, &guestProto.WidgetGuestParams{
				FirstName:   guest.FirstName,
				LastName:    guest.LastName,
				Email:       guest.Email,
				PhoneNumber: guest.PhoneNumber,
				Language:    guest.Language,
				BirthDate:   guest.BirthDate,
				Primary:     guest.Primary,
			})
		}

		proto, err := s.ReservationWidgetClient.CreateWidgetReservation(ctx, &guestProto.CreateWidgetReservationRequest{
			Guests:            guests,
			SeatingAreaId:     input.SeatingAreaId,
			BranchId:          input.BranchId,
			GuestsNumber:      input.GuestsNumber,
			Date:              input.Date,
			Time:              input.Time,
			ReservedVia:       input.ReservedVia,
			SpecialOccasionId: input.SpecialOccasionId,
			Note:              input.Note,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, openapi.CreateWidgetReservationResponse{
			Code:    proto.GetCode(),
			Message: proto.GetMessage(),
		})
	}
}
