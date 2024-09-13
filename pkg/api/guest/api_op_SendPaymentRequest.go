package guest

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) SendPaymentRequest() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var (
			input    openapi.PaymentRequest
			items    []*guestProto.PaymentItem
			delivery string
			contacts string
		)

		common.ResolveRequestBody(c, &input)

		for _, contact := range input.Contacts {
			contacts += fmt.Sprintf("%v", contact)
			if contact == input.Contacts[len(input.Contacts)-1] {
				contacts += ","
			}
		}

		for _, d := range input.Delivery {
			if d == input.Delivery[len(input.Delivery)-1] {
				delivery += strings.ToLower(d)
			} else {
				delivery += strings.ToLower(d) + ","
			}
		}

		for _, item := range input.Items {
			items = append(items, &guestProto.PaymentItem{
				Id:       item.Id,
				Name:     item.Name,
				Price:    item.Price,
				Quantity: item.Quantity,
			})
		}

		proto, err := s.PaymentClient.SendPaymentRequest(ctx, &guestProto.PaymentRequest{
			ReservationId:    common.ConvertStringToInt(c.Param("reservationId")),
			Items:            items,
			Delivery:         delivery,
			Contacts:         contacts,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, buildPaymentResponse(proto))
	}
}
