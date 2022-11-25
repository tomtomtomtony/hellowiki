package utils

import (
	"hellowiki/common"
	"strconv"
)

func ConstructCategoryNameId(categoryName string, categoryId uint) string {
	return categoryName + common.UNDER_SCORE + strconv.Itoa(int(categoryId))
}
