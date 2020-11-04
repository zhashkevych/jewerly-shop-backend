package handler

import (
	"github.com/gin-gonic/gin"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
	"gopkg.in/guregu/null.v3"
	"strconv"
)

const (
	defaultLimit  = 20
	defaultOffset = 0
)

func getProductFilters(c *gin.Context) jewerly.GetAllProductsFilters {
	filters := jewerly.GetAllProductsFilters{
		Language: jewerly.GetLanguageFromQuery(c.Query("language")),
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 {
		filters.Limit = defaultLimit
	} else {
		filters.Limit = limit
	}

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil || offset < 0{
		filters.Offset = defaultOffset
	} else {
		filters.Offset = offset
	}

	categoryId, err := strconv.Atoi(c.Query("category"))
	if err != nil {
		return filters
	}

	if err := jewerly.Category(categoryId).Validate(); err != nil {
		return filters
	}

	filters.CategoryId = null.NewInt(int64(categoryId), true)

	return filters
}

func getOrderFilters(c *gin.Context) jewerly.GetAllOrdersFilters {
	var filters jewerly.GetAllOrdersFilters

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0{
		filters.Limit = defaultLimit
	} else {
		filters.Limit = limit
	}

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil || offset < 0{
		filters.Offset = defaultOffset
	} else {
		filters.Offset = offset
	}
	return filters
}
