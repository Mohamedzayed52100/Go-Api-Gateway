package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *SUser) RequestResetPassword() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.RequestResetPasswordRequest
		common.ResolveRequestBody(c, &input)

		proto, err := s.UserClient.RequestResetPassword(ctx, &userProto.RequestResetPasswordRequest{
			Email:       input.Email,
			PhoneNumber: input.PhoneNumber,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, openapi.RequestResetPasswordResponse{
			Token: proto.GetToken(),
		})
	}
}
