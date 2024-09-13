package customerservice

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	cssProto "github.com/goplaceapp/goplace-customer-service/api/v1"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SCustomerService) GetTopEmployees() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		limit := c.Query("limit")
		if limit == "" {
			limit = "4"
		}

		proto, err := s.CSSEmployeeClient.GetTopEmployees(ctx, &cssProto.GetTopEmployeesRequest{
			Limit: common.ConvertStringToInt(limit),
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildTopEmployeesResponse(proto.GetEmployees()))
	}
}
