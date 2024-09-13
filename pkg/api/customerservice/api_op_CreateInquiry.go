package customerservice

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	cssProto "github.com/goplaceapp/goplace-customer-service/api/v1"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SCustomerService) CreateInquiry() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.CreateInquiryRequest
		common.ResolveRequestBody(c, &input)

		proto, err := s.CSSInquiryClient.CreateInquiry(ctx, &cssProto.CreateInquiryRequest{
			Params: &cssProto.InquiryParams{
				ReservationId:    input.ReservationId,
				GuestId:          input.GuestId,
				TypeId:           input.TypeId,
				Description:      input.Description,
				CreationDuration: input.CreationDuration,
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, buildInquiryResponse(proto.GetResult()))
	}
}
