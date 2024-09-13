package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *SUser) UpdatePinCode() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.UpdatePinCodeRequest
		common.ResolveRequestBody(c, &input)

		proto, err := s.UserClient.UpdatePinCode(ctx, &userProto.UpdatePinCodeRequest{
			OldPinCode:        input.OldPinCode,
			NewPinCode:        input.NewPinCode,
			ConfirmNewPinCode: input.ConfirmNewPinCode,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, openapi.UpdatePinCodeResponse{
			Code:    proto.GetCode(),
			Message: proto.GetMessage(),
		})
	}
}
