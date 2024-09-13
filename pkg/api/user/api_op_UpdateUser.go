package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *SUser) UpdateUser() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input UpdateUserForm
		if err := c.Bind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var fileName string
		avatar, _, err := c.Request.FormFile("avatar")
		if err == nil {
			fileName, err = common.UploadFileToS3(avatar, "avatars", []string{"jpg", "jpeg", "png"})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			defer avatar.Close()
		}

		userId := common.ConvertStringToInt(c.Param("userId"))
		proto, err := s.UserClient.UpdateUser(ctx, &userProto.UpdateUserRequest{
			Params: &userProto.UserParams{
				Id:          userId,
				EmployeeId:  input.EmployeeId,
				FirstName:   input.FirstName,
				LastName:    input.LastName,
				Email:       input.Email,
				PhoneNumber: input.PhoneNumber,
				Role:        common.ConvertStringToInt(input.Role),
				Avatar:      fileName,
				Birthdate:   input.Birthdate,
				JoinedAt:    input.JoinedAt,
				BranchIds:   common.ParseQueryCriteriaToIntArray(input.BranchIds),
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildUserResponse(proto.GetResult()))
	}
}
