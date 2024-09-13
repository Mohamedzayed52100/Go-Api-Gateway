package user

import (
	"context"
	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
)

func (s *SUser) GetAuthenticatedUser() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		proto, err := s.UserClient.GetAuthenticatedUser(ctx, &emptypb.Empty{})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildAuthenticatedUserResponse(proto.GetResult()))
	}
}
