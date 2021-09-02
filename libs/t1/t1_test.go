package t1

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/suhanyujie/go-utils/helper/copyer"
	"github.com/suhanyujie/go-utils/helper/mymap"
	"github.com/suhanyujie/go-utils/helper/mystring"
	"github.com/suhanyujie/go-utils/helper/slicex"
	"github.com/suhanyujie/go-utils/libs/jsonx"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"math"
	"strconv"
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
	suffix := mystring.Substr(op, strings.LastIndex(op, ".")+1, len(op))
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
	a := []int64{1145, 2, 3, 24335}
	j1 := jsonx.ToJsonIgnoreErr(a)
	t.Log(j1)
	t.Log("---end...")
}

func TestInterface1(t *testing.T) {
	l1 := []interface{}{
		"title",
		"priority",
		"status",
	}
	input := "status"
	has, err := slicex.Contain(l1, input)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(has)
}

// 浮点数，非四舍五入的转换
func TestFormat1(t *testing.T) {
	fNum := float64(216981.236)
	accuracyNum := 2
	coefficientNum := 100

	floatStr := fmt.Sprintf("%."+strconv.Itoa(accuracyNum)+"f", fNum)
	inst, _ := strconv.ParseFloat(floatStr, 64)
	t.Log(inst)
	// 千分位，非四舍五入
	fNum, _ = decimal.NewFromFloat(float64(coefficientNum)).Mul(decimal.NewFromFloat(fNum)).Float64()
	fNum = math.Floor(fNum)
	f1, _ := decimal.NewFromFloat(fNum).Div(decimal.NewFromFloat(float64(coefficientNum))).RoundBank(int32(accuracyNum)).Float64()
	t.Log(f1)

	// 打印出千分位的值 https://www.cnblogs.com/DillGao/p/8986602.html
	p := message.NewPrinter(language.English)
	renderedValStr := p.Sprintf("%."+strconv.Itoa(accuracyNum)+"f", f1)
	fmt.Printf("固定位数： %v, 千分位： %v\n", f1, renderedValStr)

	// 10.2 非四舍五入
	fNum = 10.2
	fNum, _ = decimal.NewFromFloat(float64(coefficientNum)).Mul(decimal.NewFromFloat(fNum)).Float64()
	fNum = math.Floor(fNum)
	f1, _ = decimal.NewFromFloat(fNum).Div(decimal.NewFromFloat(float64(coefficientNum))).RoundBank(int32(accuracyNum)).Float64()
	fmt.Printf("%v, %v", fNum, f1)
}

// 空切片和 nil
func TestNilSlice1(t *testing.T) {
	s1 := make([]int, 0)
	if len(s1) != 0 && s1 == nil {
		t.Error("error for s1")
		return
	}
	s1 = nil
	if len(s1) != 0 {
		t.Error("error for s1-1")
		return
	}
}

type OneOp1 struct {
	Id    interface{} `json:"id"`
	Color string      `json:"color"`
	Value string      `json:"value"`
}

func TestUnmarshal1(t *testing.T) {
	cellVal := float64(3)
	json1 := `[{"color": "", "id": 3, "value": "老板1"},{"color": "", "id": "123", "value": "老板2"}]`
	obj1 := make([]OneOp1, 0)
	jsonx.FromJson(json1, &obj1)
	map1 := GetOpList(obj1)

	if res1, ok := map1[cellVal]; ok {
		t.Log(jsonx.ToJsonIgnoreErr(res1))
		return
	}
	fmt.Printf("res: %v\n", obj1[0].Id.(float64) == 3)
	fmt.Println(jsonx.ToJsonIgnoreErr(obj1))
}

func GetOpList(ops []OneOp1) map[interface{}]OneOp1 {
	map1 := make(map[interface{}]OneOp1, 0)
	for _, op := range ops {
		map1[op.Id] = op
	}
	return map1
}

func TestCaseCamelCopy(t *testing.T) {
	res := mymap.CaseCamelCopy(map[string]interface{}{
		"issue_status": 1,
	})
	t.Log(res)
}

type Issue struct {
	Title  string
	LcData map[string]interface{}
}

// 当 map 为 nil 时，直接赋值是会报 panic 的
func TestNilMap1(t *testing.T) {
	issue1 := Issue{
		Title: "t1",
	}
	if issue1.LcData == nil {
		issue1.LcData = make(map[string]interface{}, 0)
	}
	issue1.LcData["title"] = issue1.Title
}

func TestMapAndAssert(t *testing.T) {
	m1 := make(map[string]interface{}, 0)
	m1["orgIds"] = []int64{1102}
	// 如果 key 不存在，也不会发生 panic
	tmpData, isOk := m1["orgIds"].([]int64)
	if isOk {
		t.Log(jsonx.ToJsonIgnoreErr(tmpData))
	} else {
		t.Log("no key data")
	}
	// 如果 key 不存在，获取到的就是 `nil`
	d2 := m1["userName"]
	d3, isOk := m1["userName"]
	t.Log(jsonx.ToJsonIgnoreErr(d2))
	t.Log(d3, isOk)
}

func TestSubStr1(t *testing.T) {
	prev := ""
	opCode := "Permission.Pro.View-ManagePrivate"
	opCode = "Permission.Pro.View.ManagePrivate"
	opCode = "Permission.Pro.Issue.4-Modify"
	opCode = "Permission.Pro.Issue.4.Modify"
	if strings.IndexAny(opCode, "-") != -1 {
		info := strings.Split(opCode, "-")
		if len(info) > 0 {
			prev = info[0]
		}
	} else {
		prev = mystring.Substr(opCode, 0, strings.LastIndex(opCode, "."))
	}
	t.Log(prev)
}

type TestObjConvert1Type1 struct {
	Num1 int `json:"int64"`
}

type TestObjConvert1Type2 struct {
	Num1 int `json:"int"`
}

// 不同类型的 field，是无法正常转换的。如 int 和 int64 是不同类型。
func TestObjConvert1(t *testing.T) {
	o1 := TestObjConvert1Type1{Num1: 1}
	dstO1 := TestObjConvert1Type2{}
	copyer.Copy(o1, &dstO1)
	t.Log(jsonx.ToJsonIgnoreErr(dstO1))
}

func TestIntJson1(t *testing.T) {
	json1 := "[997639]"
	arr := make([]int64, 0)
	err := jsonx.FromJson(json1, &arr)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(jsonx.ToJsonIgnoreErr(arr))
}

func TestJsonToInterface1(t *testing.T) {
	map1 := map[string]interface{}{
		"userAge": 18,
	}
	json1 := jsonx.ToJsonIgnoreErr(map1)
	map2 := make(map[string]interface{}, 0)
	jsonx.FromJson(json1, &map2)
	age, ok1 := map2["userAge"]
	// age1 := age.(int) // bad
	age1 := int(age.(float64)) // ok
	t.Log(age, ok1, age1)
}
