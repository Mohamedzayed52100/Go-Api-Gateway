package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *SUser) GetUserByID() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		userId := common.ConvertStringToInt(c.Param("userId"))
		proto, err := s.UserClient.GetUserByID(ctx, &userProto.GetUserByIDRequest{
			Id: userId,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildUserResponse(proto.GetResult()))
	}
}
