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

func TestFromJsonWithNumber(t *testing.T) {
	json1 := `[9819899]`
	arr := make([]int64, 0)
	err := FromJson(json1, &arr)
	if err != nil {
		t.Error(err)
		return
	}
	resStr := ToJsonIgnoreErr(arr)
	if json1 != resStr {
		t.Error("error convert 001")
		return
	}
	t.Log(resStr)
}
