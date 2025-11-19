package redisx

import (
	"testing"
)

func TestNewClient1(t *testing.T) {
	// initial redis client todo
	Init(nil)
	cli1 := GetRedisClient()
	t1, err := cli1.Time().Result()
	if err != nil {
		t.Logf("[TestNewClient1] err: %v", err)
		return
	}
	t.Logf("res: %v", t1)
}

func TestLock1(t *testing.T) {
	key := "game:farm:farm:uid:1"
	// set conf todo
	isOk := GetLock(key)
	defer UnLock(key)
	if isOk {
		t.Logf("get lock ok")
	} else {
		t.Logf("get lock fail")
	}
}

func TestSetKey(t *testing.T) {
	key := "testSet"
	val := "11111"
	cli := GetRedisClient()
	cli.Set(key, val, 0)
	value := cli.Get(key)
	t.Logf("value = %s", value.Val())
}
