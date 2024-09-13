package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *SUser) CheckPinCode() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.PinCodeRequest
		common.ResolveRequestBody(c, &input)
		role := c.Query("role")

		proto, err := s.UserClient.CheckPinCode(ctx, &userProto.PinCodeRequest{
			PinCode: input.PinCode,
			Role:    role,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildPinCodeResponse(proto.GetResult()))
	}
}
