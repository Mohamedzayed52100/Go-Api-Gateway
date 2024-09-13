package guest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *SGuest) CreateGuest() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var input openapi.CreateGuestRequest
		common.ResolveRequestBody(c, &input)
		var (
			tags      []*guestProto.TagParams
			email     string
			language  string
			birthdate string
		)
		for _, tag := range input.Tags {
			tags = append(tags, &guestProto.TagParams{
				Id:         tag.Id,
				CategoryId: tag.CategoryId,
			})
		}

		if input.Email != nil {
			email = *input.Email
		}
		if input.Language != nil {
			language = *input.Language
		}

		if input.BirthDate != nil {
			birthdate = *input.BirthDate
		}

		if input.Gender == "" {
			input.Gender = "male"
		}

		proto, err := s.GuestClient.CreateGuest(ctx, &guestProto.CreateGuestRequest{
			Params: &guestProto.GuestParams{
				FirstName:   input.FirstName,
				LastName:    input.LastName,
				Email:       email,
				PhoneNumber: input.PhoneNumber,
				Language:    language,
				BirthDate:   birthdate,
				Tags:        tags,
				Address:     input.Address,
				Gender:      input.Gender,
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusCreated, buildGuestAPIResponse(proto.GetResult()))
	}
}
