package temptools

import (
	"strings"

	"github.com/valyala/fasttemplate"
)

// ReplaceTemplate 模版替换
func ReplaceTemplate(temp string, v map[string]interface{}) string {
	t := fasttemplate.New(temp, "{{", "}}")
	return t.ExecuteString(v)
}

// ToUpperFistring 将字符串的首字符大写
func ToUpperFistring(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}
