package customerservice

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	cssProto "github.com/goplaceapp/goplace-customer-service/api/v1"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SCustomerService) GetInquiryCommentByID() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		inquiryId := common.ConvertStringToInt(c.Param("inquiryId"))
		commentId := common.ConvertStringToInt(c.Param("commentId"))

		proto, err := s.CSSInquiryCommentClient.GetInquiryCommentByID(ctx, &cssProto.GetInquiryCommentByIDRequest{
			InquiryId: inquiryId,
			CommentId: commentId,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, buildInquiryCommentResponse(proto.GetResult()))
	}
}
