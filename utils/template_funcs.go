package utils

import (
	"github.com/beego/beego/v2/server/web"
	"strings"
)

func Multiply(a, b int) int {
	return a * b
}

// SanitizeID replaces non-alphanumeric characters with underscores
func SanitizeID(id string) string {
	return strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return '_'
	}, id)
}

// RegisterTemplateFuncs registers all custom template functions with Beego
func RegisterTemplateFuncs() {
	web.AddFuncMap("mul", Multiply)
	web.AddFuncMap("sanitizeID", SanitizeID)
}
