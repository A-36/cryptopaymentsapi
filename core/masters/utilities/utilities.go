package utils

import (
	"fmt"
	"regexp"
)

func SafeString(str string) string {
	reg, err := regexp.Compile("[^A-Za-z0-9]+")
	if err != nil {
		fmt.Println(err)
	}
	newStr := reg.ReplaceAllString(str, "")
	return newStr
}
