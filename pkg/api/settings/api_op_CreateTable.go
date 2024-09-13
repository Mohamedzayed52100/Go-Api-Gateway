package settings

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
)

func (s *SSettings) CreateTable() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.CreateTableRequest
		common.ResolveRequestBody(c, &input)
		combinedTables := ""
		for _, v := range input.CombinedTables {
			if v == input.CombinedTables[len(input.CombinedTables)-1] {
				combinedTables += common.ConvertIntToString(v)
			} else {
				combinedTables += common.ConvertIntToString(v) + ","
			}
		}
		proto, err := s.TableClient.CreateTable(ctx, &settingsProto.CreateTableRequest{
			Params: &settingsProto.TableParams{
				SeatingAreaId:  input.SeatingAreaId,
				TableNumber:    input.TableNumber,
				PosNumber:      input.PosNumber,
				MinPartySize:   input.MinPartySize,
				MaxPartySize:   input.MaxPartySize,
				CombinedTables: combinedTables,
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, buildTableResponse(proto.GetResult()))
	}
}
