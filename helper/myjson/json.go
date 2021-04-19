package myjson

import "encoding/json"

// 将对象转换为 json 字符串
func ToJsonIgnore(obj interface{}) string {
	v, _ := json.Marshal(obj)
	return string(v)
}
