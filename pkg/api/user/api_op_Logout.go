package user

import (
	"context"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *SUser) Logout() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		_, err := s.UserClient.Logout(ctx, &emptypb.Empty{})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(200, gin.H{
			"message": "Successfully logout",
		})
	}
}
