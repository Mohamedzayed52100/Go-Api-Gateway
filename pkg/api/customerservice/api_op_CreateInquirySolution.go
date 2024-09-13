package customerservice

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	cssProto "github.com/goplaceapp/goplace-customer-service/api/v1"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SCustomerService) CreateInquirySolution() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.AddInquirySolutionRequest
		common.ResolveRequestBody(c, &input)
		inquiryId := common.ConvertStringToInt(c.Param("inquiryId"))

		proto, err := s.CSSInquiryClient.CreateInquirySolution(ctx, &cssProto.CreateInquirySolutionRequest{
			Params: &cssProto.InquirySolutionParams{
				InquiryId: inquiryId,
				Solution:  input.Solution,
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, buildInquirySolutionResponse(proto.GetResult()))
	}
}
