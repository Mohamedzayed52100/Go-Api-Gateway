package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *SUser) RequestDemo() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.DemoRequest
		common.ResolveRequestBody(c, &input)

		proto, err := s.TenantClient.RequestDemo(ctx, &userProto.DemoRequest{
			Params: &userProto.DemoParams{
				Name:           input.Name,
				Email:          input.Email,
				PhoneNumber:    input.PhoneNumber,
				Country:        input.Country,
				RestaurantName: input.RestaurantName,
				BranchesNo:     int32(input.BranchesNo),
				FirstTimeCrm:   input.FirstTimeCrm,
				SystemName:     input.SystemName,
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, openapi.DemoResponse{
			Result: proto.GetResult(),
		})
	}
}
