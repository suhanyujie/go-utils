package t1

import (
	"errors"
	"github.com/suhanyujie/go-utils/helper/mystring"
	"github.com/suhanyujie/go-utils/libs/jsonx"
	"strings"
	"testing"
)

func TestLastIndex(t *testing.T) {
	op := "xxx.xxxConfig.Create"
	prev := mystring.Substr(op, 0, strings.LastIndex(op, "."))
	if prev != "xxx.xxxConfig" {
		t.Error(errors.New("error 001"))
		return
	}
	t.Log("end...")
}

func TestLastIndex1(t *testing.T) {
	op := "xxx.xxxConfig.Create"
	suffix := mystring.Substr(op, strings.LastIndex(op, ".") + 1, len(op))
	if suffix != "Create" {
		t.Error(errors.New("error 001"))
		return
	}
	t.Log("end...")
}

type UserIdsObj struct {
	UserIds []int64 `json:"userIds"`
}

func TestIfJson(t *testing.T) {
	a := []int64{1145,2,3, 24335}
	j1 := jsonx.ToJsonIgnore(a)
	t.Log(j1)
	t.Log("---end...")
}
