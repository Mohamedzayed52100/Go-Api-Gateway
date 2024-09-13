package user

import (
	"context"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SUser) SendMessage() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.SendMessageRequest
		common.ResolveRequestBody(c, &input)

		// proto, err := s.UserClient.SendMessage(ctx, &userProto.SendMessageRequest{
		// 	To:      input.To,
		// 	Message: input.Message,
		// })
		// if err != nil {
		// 	common.HandleGrpcError(c, err)
		// 	return
		// }

		// c.JSON(http.StatusOK, openapi.SendMessageResponse{
		// 	Code:    proto.Code,
		// 	Message: proto.Message,
		// })
	}
}
