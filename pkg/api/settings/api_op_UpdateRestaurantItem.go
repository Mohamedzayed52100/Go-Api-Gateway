package settings

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/goplaceapp/goplace-common/pkg/meta"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SSettings) UpdateRestaurantItem() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input UpdateRestaurantItemDto

		if err := c.Bind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var fileName string
		image, header, err := c.Request.FormFile("image")
		if err == nil {
			// Check if the file size is greater than 2 MB
			const maxFileSize = 2 << 20
			if header.Size > maxFileSize {
				c.JSON(http.StatusBadRequest, gin.H{"error": "File size is greater than 2 MB"})
				return
			}

			// Upload the file to S3
			folderName := fmt.Sprintf("clients/%s/menu-items/%d", ctx.Value(meta.TenantIDContextKey.String()).(string), time.Now().Unix())
			fileName, err = common.UploadFileToS3(image, folderName, []string{"jpg", "jpeg", "png"})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			defer image.Close()
		}

		proto, err := s.RestaurantItemClient.UpdateRestaurantItem(ctx, &settingsProto.UpdateRestaurantItemRequest{
			Params: &settingsProto.RestaurantItemParams{
				Id:    common.ConvertStringToInt(c.Param("itemId")),
				Name:  input.Name,
				Price: input.Price,
				Image: fileName,
				Icon:  input.Icon,
				Code:  input.Code,
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildRestaurantItemResponse(proto.GetResult()))
	}
}
