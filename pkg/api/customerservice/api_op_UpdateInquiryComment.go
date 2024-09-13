package customerservice

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	cssProto "github.com/goplaceapp/goplace-customer-service/api/v1"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SCustomerService) UpdateInquiryComment() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.UpdateInquiryCommentRequest
		common.ResolveRequestBody(c, &input)
		commentId := common.ConvertStringToInt(c.Param("commentId"))
		inquiryId := common.ConvertStringToInt(c.Param("inquiryId"))

		proto, err := s.CSSInquiryCommentClient.UpdateInquiryComment(ctx, &cssProto.UpdateInquiryCommentRequest{
			Params: &cssProto.InquiryCommentParams{
				Id:        commentId,
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
