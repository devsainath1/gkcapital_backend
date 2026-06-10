package utils

import (
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	nameRegex  = regexp.MustCompile(`^[a-zA-Z\s]+$`)
	phoneRegex = regexp.MustCompile(`^[0-9]{10}$`)
	cityRegex  = regexp.MustCompile(`^[a-zA-Z\s]+$`)
)

func ValidateName(name string) bool {
	return nameRegex.MatchString(name)
}

func ValidatePhone(phone string) bool {
	return phoneRegex.MatchString(phone)
}

func ValidateCity(city string) bool {
	return cityRegex.MatchString(city)
}


func GetPagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return page, pageSize
}

func GetOffset(page, pageSize int) int {
	return (page - 1) * pageSize
}
