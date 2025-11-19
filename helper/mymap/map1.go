package mymap

import "github.com/suhanyujie/go_utils/helper/mystring"

// CaseCamelCopy 将 map 的下划线风格的键转换为小驼峰风格的键
func CaseCamelCopy(data map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{}, 0)
	for s, i := range data {
		key := mystring.LcFirst(mystring.Case2Camel(s))
		res[key] = i
	}

	return res
}
