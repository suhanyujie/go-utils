package jsonx

import (
	"testing"
)

func TestToJson(t *testing.T) {
	type Stu struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	json := ToJsonIgnoreErr(Stu{
		Name: "李LiuDeHua",
		Age:  21,
	})
	t.Log(json)
}
