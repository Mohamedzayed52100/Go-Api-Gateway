package guest

import (
	"context"
	"net/http"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"

	"github.com/gin-gonic/gin"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

func (s *SGuest) UpdateGuest() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var (
			input          openapi.UpdateGuestRequest
			tags           []*guestProto.TagParams
			email          string
			emptyEmail     bool
			language       string
			emptyLanguage  bool
			birthdate      string
			emptyBirthdate bool
			emptyTags      bool
		)

		id := c.Param("guestId")
		common.ResolveRequestBody(c, &input)

		if input.Tags == nil {
			emptyTags = true
		} else {
			for _, tag := range input.Tags {
				tags = append(tags, &guestProto.TagParams{
					Id:         tag.Id,
					CategoryId: tag.CategoryId,
				})
			}
		}

		if input.Email != nil {
			email = *input.Email
		} else {
			emptyEmail = true
		}

		if input.Language != nil {
			language = *input.Language
		} else {
			emptyLanguage = true
		}

		if input.BirthDate != nil {
			birthdate = *input.BirthDate
		} else {
			emptyBirthdate = true
		}

		proto, err := s.GuestClient.UpdateGuest(ctx, &guestProto.UpdateGuestRequest{
			Params: &guestProto.GuestParams{
				Id:             common.ConvertStringToInt(id),
				FirstName:      input.FirstName,
				LastName:       input.LastName,
				Email:          email,
				PhoneNumber:    input.PhoneNumber,
				Language:       language,
				BirthDate:      birthdate,
				Tags:           tags,
				EmptyEmail:     emptyEmail,
				EmptyLanguage:  emptyLanguage,
				EmptyBirthdate: emptyBirthdate,
				EmptyTags:      emptyTags,
				Address:        input.Address,
				Gender:         input.Gender,
			},
		})
		if err != nil {
			common.HandleGrpcError(c, err)
			return
		}

		c.JSON(http.StatusOK, buildGuestAPIResponse(proto.GetResult()))
	}
}
