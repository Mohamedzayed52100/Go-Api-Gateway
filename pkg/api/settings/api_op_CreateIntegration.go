package settings

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/goplaceapp/goplace-common/pkg/logger"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
	"net/http"
)

func (s *SSettings) CreateIntegration() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.CreateIntegrationRequest
		common.ResolveRequestBody(c, &input)

		credentialsJson, err := json.Marshal(input.Credentials)
		if err != nil {
			logger.Default().Errorf("Error marshalling credentials: %v", err)
			return
		}

		proto, err := s.SettingsClient.CreateIntegration(ctx, &settingsProto.CreateIntegrationRequest{
			Params: &settingsProto.IntegrationParams{
				SystemName:     input.SystemName,
				SystemType:     input.SystemType,
				BaseURL:        input.BaseURL,
				CredentialType: input.CredentialType,
				Credentials:    string(credentialsJson),
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, buildIntegrationResponse(proto.GetResult()))
	}
}
