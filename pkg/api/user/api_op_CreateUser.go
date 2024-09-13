package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *SUser) CreateUser() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input CreateUserForm
		var fileName string

		if err := c.Bind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		avatar, _, err := c.Request.FormFile("avatar")
		if err == nil {
			fileName, err = common.UploadFileToS3(avatar, "avatars", []string{"jpg", "jpeg", "png"})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			defer avatar.Close()
		}

		proto, err := s.UserClient.CreateUser(ctx, &userProto.CreateUserRequest{
			Params: &userProto.UserParams{
				EmployeeId:  input.EmployeeId,
				FirstName:   input.FirstName,
				LastName:    input.LastName,
				Email:       input.Email,
				PhoneNumber: input.PhoneNumber,
				Role:        common.ConvertStringToInt(input.Role),
				Avatar:      fileName,
				JoinedAt:    input.JoinedAt,
				Birthdate:   input.Birthdate,
				BranchIds:   common.ParseQueryCriteriaToIntArray(input.BranchIds),
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, buildUserResponse(proto.GetResult()))
	}
}
