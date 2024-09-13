package user

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *SUser) Login() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.LoginRequest
		common.ResolveRequestBody(c, &input)

		proto, err := s.UserClient.Login(ctx, &userProto.LoginRequest{
			Params: &userProto.LoginParams{
				Email:        input.Email,
				Password:     input.Password,
				ClientID:     input.ClientID,
				ClientSecret: input.ClientSecret,
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		expirationTime := proto.GetResult().GetExpiresAt().AsTime()

		cookie := &http.Cookie{
			Name:     meta.AccessTokenCookieName,
			Path:     "/",
			Value:    proto.GetResult().GetAccessToken(),
			Expires:  expirationTime,
			HttpOnly: true,
		}

		http.SetCookie(c.Writer, cookie)

		setup := proto.GetResult().GetSetup()

		res := &openapi.LoginResponse{
			AccessToken: proto.GetResult().GetAccessToken(),
			ExpiresAt:   strconv.FormatInt(expirationTime.Unix(), 10),
		}

		if !setup {
			res.Setup = &setup
		} else {
			res.Setup = nil
		}

		c.JSON(http.StatusOK, res)
	}
}
