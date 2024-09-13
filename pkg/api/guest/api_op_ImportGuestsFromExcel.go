package guest

import (
	"context"
	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"net/http"
)

func (s *SGuest) ImportGuestsFromExcel() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var (
			filePath          string
			allowedExtensions = []string{"xlsx", "xls", "csv"}
		)

		file, _, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		filePath, err = common.UploadFileToS3(file, "GuestDataImports", allowedExtensions)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer file.Close()

		proto, err := s.GuestClient.ImportGuestsFromExcel(ctx, &guestProto.ImportGuestsFromExcelRequest{
			FilePath: filePath,
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, &openapi.ImportGuestsFromExcelResponse{
			Code:    proto.GetCode(),
			Message: proto.GetMessage(),
		})
	}
}
