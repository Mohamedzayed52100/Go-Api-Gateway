package customerservice

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	cssProto "github.com/goplaceapp/goplace-customer-service/api/v1"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SCustomerService) GetGuestInquiries() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		id := common.ConvertStringToInt(c.Param("guestId"))
		proto, err := s.CSSInquiryClient.GetGuestInquiries(ctx, &cssProto.GetGuestInquiriesRequest{
			Id: id,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildAllInquiriesResponse(proto.GetResult()))
	}
}
