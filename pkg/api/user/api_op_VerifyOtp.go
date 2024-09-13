package user

import (
	"context"
	"net/http"

	userProto "github.com/goplaceapp/goplace-user/api/v1"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SUser) VerifyOtp() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.VerifyOtpRequest
		common.ResolveRequestBody(c, &input)

		proto, err := s.UserClient.VerifyOtp(ctx, &userProto.VerifyOtpRequest{
			Code:  input.Code,
			Token: input.Token,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, openapi.VerifyOtpResponse{
			Code:    proto.GetCode(),
			Message: proto.GetMessage(),
		})
	}
}
