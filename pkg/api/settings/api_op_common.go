package settings

import (
	"sort"
	"strings"

	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
)

type SSettings struct {
	*common.Resources
}

type CreateRestaurantItemDto struct {
	Name  string  `form:"name" binding:"required"`
	Price float32 `form:"price" binding:"required"`
	Code  string  `form:"code" binding:"required"`
	Icon  string  `form:"icon"`
}

type UpdateRestaurantItemDto struct {
	Name  string  `form:"name"`
	Price float32 `form:"price"`
	Code  string  `form:"code"`
	Icon  string  `form:"icon"`
}

func SortItemsByCode(items []*settingsProto.RestaurantItem, criteria string) {
	sort.Slice(items, func(i, j int) bool {
		if criteria == "asc" {
			return strings.ToLower(items[i].Code) < strings.ToLower(items[j].Code)
		} else {
			return strings.ToLower(items[i].Code) > strings.ToLower(items[j].Code)
		}
	})
}

func SortItemsByName(items []*settingsProto.RestaurantItem, criteria string) {
	sort.Slice(items, func(i, j int) bool {
		if criteria == "asc" {
			return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
		} else {
			return strings.ToLower(items[i].Name) > strings.ToLower(items[j].Name)
		}
	})
}

func ListLowerCase(items []string) []string {
	result := make([]string, len(items))
	for i := range items {
		result[i] = strings.ToLower(items[i])
	}

	return result
}
