package temptools

import "github.com/valyala/fasttemplate"

// ReplaceTemplate 模版替换
func ReplaceTemplate(temp string, v map[string]interface{}) string {
	t := fasttemplate.New(temp, "{{", "}}")
	return t.ExecuteString(v)
}
