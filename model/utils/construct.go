package utils

import (
	"hellowiki/common"
	"strconv"
)

func ConstructCategoryNameId(categoryEngName string, categoryId uint) string {
	return ConstructStandardIndexName(categoryEngName, categoryId)
}

func ConstructStandardIndexName(categoryEngName string, categoryId uint) string {
	return categoryEngName + common.UNDER_SCORE + strconv.Itoa(int(categoryId))
}
