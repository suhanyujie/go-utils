package jsonx

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestCalVal1(t *testing.T) {
	val1 := 1 << 8
	t.Log(val1)
}

func TestSyncMap1(t *testing.T) {
	m1 := sync.Map{}
	m1.Store("name", "LiYuChun")
	m1.Store("age", 12)
	userName := ""
	m1.Range(func(key, value any) bool {
		valStr, ok := value.(string)
		if !ok {
			return false
		}
		userName = valStr
		return true
	})
	t.Logf("userName: %s", userName)
}

// sync map 在遍历是删除元素操作
func TestSyncMapDelete1(t *testing.T) {
	m1 := sync.Map{}
	m1.Store("name", "dengJiaXian")
	m1.Store("age", 20)
	m1.Range(func(key, value any) bool {
		if key.(string) == "age" {
			m1.Delete(key)
		}
		return true
	})
	m1.Range(func(key, value any) bool {
		fmt.Printf("key: %v, value: %v \n", key, value)
		return true
	})
}

func TestSwitch1(t *testing.T) {
	typ := 1
	switch typ {
	case 1:
		fmt.Printf("1-test1\n")
		break
		fmt.Printf("1-test2\n")
	case 2:
		fmt.Printf("2-test2\n")
	}
}

type Stu1 struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestStruct1(t *testing.T) {
	stu1 := Stu1{
		Name: "name 1001",
		Age:  21,
	}
	stu2 := stu1
	stu2.Age = 22
	t.Logf("stu1: %s, stu2： %s", ToJsonIgnoreErr(stu1), ToJsonIgnoreErr(stu2))
}

func TestTick1(t *testing.T) {
	tic := time.NewTicker(2 * time.Second)
	for i := 0; i < 100; i++ {
		select {
		case <-tic.C:
			newInterval := rand.Intn(10)
			if newInterval > 0 {
				tic.Reset(time.Duration(newInterval) * time.Second)
			}
			log.Printf("task 1-1 newInterval：%v", newInterval)
		}
	}
}

func TestRecover1(t *testing.T) {
	log.Printf("1-1")
	recover()
	log.Printf("1-2")
}

func TestFloat1(t *testing.T) {
	val := 25
	res1 := val / 2
	res2 := float32(val) / 2
	t.Logf("res1: %v, res2: %v", res1, res2)
}

func TestToJsonFormatIgnoreErr(t *testing.T) {
	json1 := `{"name":"name 1001","age":21}`
	m1 := make(map[string]interface{})
	FromJson(json1, &m1)
	res := ToJsonFormatIgnoreErr(m1)
	t.Logf("res: %s", res)
}
