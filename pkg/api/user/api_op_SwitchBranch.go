package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *SUser) SwitchBranch() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		branchId := common.ConvertStringToInt(c.Param("branchId"))
		proto, err := s.UserClient.SwitchBranch(ctx, &userProto.SwitchBranchRequest{
			BranchId: branchId,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    proto.GetCode(),
			"message": proto.GetMessage(),
		})
	}
}
