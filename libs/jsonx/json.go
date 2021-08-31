package jsonx

import (
	"encoding/json"
	"strings"
)

// ToJson 将对象转换为 json 字符串
func ToJson(obj interface{}) (string, error) {
	v, err := json.Marshal(obj)
	return string(v), err
}

// ToJsonIgnoreErr 将对象转换为 json 字符串，发生异常时会被忽略
func ToJsonIgnoreErr(obj interface{}) string {
	v, _ := json.Marshal(obj)
	return string(v)
}

// FromJson json 的反序列化
func FromJson(jsonStr string, obj interface{}) error {
	return json.Unmarshal([]byte(jsonStr), obj)
}

func FromJsonWithNumber(jsonStr string, obj interface{}) error {
	d := json.NewDecoder(strings.NewReader(jsonStr))
	d.UseNumber()
	return d.Decode(&obj)
}
