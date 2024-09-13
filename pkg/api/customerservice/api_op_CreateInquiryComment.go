package customerservice

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	cssProto "github.com/goplaceapp/goplace-customer-service/api/v1"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SCustomerService) CreateInquiryComment() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.CreateInquiryCommentRequest
		common.ResolveRequestBody(c, &input)
		inquiryId := common.ConvertStringToInt(c.Param("inquiryId"))

		proto, err := s.CSSInquiryCommentClient.CreateInquiryComment(ctx, &cssProto.CreateInquiryCommentRequest{
			Params: &cssProto.InquiryCommentParams{
				InquiryId: inquiryId,
				Comment:   input.Comment,
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, buildInquiryCommentResponse(proto.GetResult()))
	}
}
