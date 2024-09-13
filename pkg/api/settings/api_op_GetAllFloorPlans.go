package settings

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *SSettings) GetAllFloorPlans() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		proto, err := s.FloorPlanClient.GetAllFloorPlans(ctx, &emptypb.Empty{})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusAccepted, buildAllFloorPlansResponse(proto.GetResult()))
	}
}
