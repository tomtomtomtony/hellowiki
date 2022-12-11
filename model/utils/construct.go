package utils

import (
	"hellowiki/common"
)

func ConstructCategoryNameId(categoryName string, uuid string) string {
	return categoryName + common.UNDER_SCORE + uuid
}
