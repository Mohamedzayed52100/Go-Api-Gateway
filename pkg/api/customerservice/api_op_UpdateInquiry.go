package customerservice

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	cssProto "github.com/goplaceapp/goplace-customer-service/api/v1"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SCustomerService) UpdateInquiry() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.UpdateInquiryRequest
		common.ResolveRequestBody(c, &input)
		id := common.ConvertStringToInt(c.Param("inquiryId"))

		proto, err := s.CSSInquiryClient.UpdateInquiry(ctx, &cssProto.UpdateInquiryRequest{
			Params: &cssProto.InquiryParams{
				Id:            id,
				ReservationId: input.ReservationId,
				Status:        input.Status,
				SolutionId:    input.SolutionId,
				TypeId:        input.TypeId,
				Description:   input.Description,
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, buildInquiryResponse(proto.GetResult()))
	}
}
