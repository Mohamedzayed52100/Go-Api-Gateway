package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *SUser) ChangePassword() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.ChangePasswordRequest
		common.ResolveRequestBody(c, &input)

		proto, err := s.UserClient.ChangePassword(ctx, &userProto.ChangePasswordRequest{
			OldPassword:        input.OldPassword,
			NewPassword:        input.NewPassword,
			ConfirmNewPassword: input.ConfirmNewPassword,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, openapi.ChangePasswordResponse{
			Code:    proto.GetCode(),
			Message: proto.GetMessage(),
		})
	}
}
