package timer

import (
	"fmt"
	"reflect"
	"time"
)

// ITimer 定时器的使用接口
type ITimer interface {
	RegTimer(func(), time.Duration) int
	RegTimerAt(func(), int, int, int, bool) int
	UnregTimer(int) error
	SuspendTimer(int) error
	ResumeTimer(int) error
	AddDelayCall(func(), time.Duration) int
	RemoveDelayCall(int) error

	RegTimerByObj(interface{}, func(), time.Duration)
	UnregTimerByObj(interface{}) error
	SuspendTimerByObj(interface{}) error
	ResumeTimerByObj(interface{}) error
	AddDelayCallByObj(interface{}, func(), time.Duration)
	RemoveDelayCallByObj(interface{}) error
}

// NewTimer 创建新的计时调度器
func NewTimer() *Timer {
	t := &Timer{}
	t.handleSeed = 1
	t.callbacks = make(map[int]*funcInfo)
	t.cbObjMap = make(map[reflect.Value][]int)
	return t
}

// Timer 实际上的Timer用法
type Timer struct {
	callbacks map[int]*funcInfo
	cbObjMap  map[reflect.Value][]int

	handleSeed int
}

type funcInfo struct {
	suspended bool
	oneTime   bool
	isAdjust  bool
	interval  time.Duration
	lastCall  time.Time
	proc      func()
}

// RegTimer 注册一个定时器
// callback为定时函数 , intervalSec 为间隔调用秒数,返回值为定时间器句柄
// 注册成功后，会定时间调用 callback
func (t *Timer) RegTimer(callback func(), intervalSec time.Duration) int {
	return t.doReg(callback, intervalSec, false)
}

// UnregTimer 反注册Timer，传入计时器句柄
func (t *Timer) UnregTimer(timerHandle int) error {
	if _, ok := t.callbacks[timerHandle]; !ok {
		return fmt.Errorf("找不到Timer: %d", timerHandle)
	}
	delete(t.callbacks, timerHandle)
	return nil
}

// SuspendTimer 挂起计时器，传入的是计时器句柄
func (t *Timer) SuspendTimer(timerHandle int) error {
	if _, ok := t.callbacks[timerHandle]; !ok {
		return fmt.Errorf("找不到Timer: %d", timerHandle)
	}
	t.callbacks[timerHandle].suspended = true
	return nil
}

// ResumeTimer 重新恢复运行计时器
func (t *Timer) ResumeTimer(timerHandle int) error {
	if _, ok := t.callbacks[timerHandle]; !ok {
		return fmt.Errorf("找不到Timer: %d", timerHandle)
	}
	t.callbacks[timerHandle].suspended = false
	return nil
}

// AddDelayCall 添加一个延时调用
func (t *Timer) AddDelayCall(callback func(), delaySec time.Duration) int {
	return t.doReg(callback, delaySec, true)
}

// RemoveDelayCall 删除一个延时调用
func (t *Timer) RemoveDelayCall(timerHandle int) error {
	return t.UnregTimer(timerHandle)
}

// GetLeftNano 获取某个定时器还剩多少纳秒结束
func (t *Timer) GetLeftNano(timerHandle int) int64 {
	handler, ok := t.callbacks[timerHandle]
	if !ok {
		fmt.Errorf("GetLeftNano找不到Timer: %d", timerHandle)
		return 0
	}

	curTime := time.Now()
	curUnixNano := curTime.UnixNano()
	elapsedNano := curUnixNano - handler.lastCall.UnixNano()
	if int64(handler.interval) > elapsedNano {
		return int64(handler.interval) - elapsedNano
	} else {
		return 0
	}
}

// RegTimerByObj 根据对象注册
func (t *Timer) RegTimerByObj(objInst interface{}, callback func(), intervalSec time.Duration) {
	v := reflect.ValueOf(objInst)
	h := t.doReg(callback, intervalSec, false)
	t.cbObjMap[v] = append(t.cbObjMap[v], h)
}

// UnregTimerByObj 去除所有该对象注册的timer
func (t *Timer) UnregTimerByObj(objInst interface{}) error {
	v := reflect.ValueOf(objInst)
	if hList, ok := t.cbObjMap[v]; ok {
		for _, h := range hList {
			if err := t.UnregTimer(h); err != nil {
				continue
			}
		}

		t.cbObjMap[v] = nil
		delete(t.cbObjMap, v)
	}

	return nil
}

// SuspendTimerByObj 暂停所有该对象注册的timer
func (t *Timer) SuspendTimerByObj(objInst interface{}) error {
	v := reflect.ValueOf(objInst)
	if hList, ok := t.cbObjMap[v]; ok {
		for _, h := range hList {
			if err := t.SuspendTimer(h); err != nil {
				return err
			}
		}
	}

	return nil
}

// ResumeTimerByObj 重新恢复运行计时器
func (t *Timer) ResumeTimerByObj(objInst interface{}) error {
	v := reflect.ValueOf(objInst)
	if hList, ok := t.cbObjMap[v]; ok {
		for _, h := range hList {
			if err := t.ResumeTimer(h); err != nil {
				return err
			}
		}
	}

	return nil
}

// AddDelayCallByObj 添加一个延时调用
func (t *Timer) AddDelayCallByObj(objInst interface{}, callback func(), intervalSec time.Duration) {
	v := reflect.ValueOf(objInst)
	h := t.doReg(callback, intervalSec, true)
	t.cbObjMap[v] = append(t.cbObjMap[v], h)
}

// RemoveDelayCallByObj 删除一个延时调用
func (t *Timer) RemoveDelayCallByObj(objInst interface{}) error {
	return t.UnregTimerByObj(objInst)
}

// Loop 计时器驱动函数
func (t *Timer) Loop() {
	curTime := time.Now()
	curUnixNano := curTime.UnixNano()
	for handle, f := range t.callbacks {
		if f.suspended == true || curUnixNano-f.lastCall.UnixNano() < int64(f.interval) {
			continue
		}
		if f.oneTime == true {
			t.UnregTimer(handle)
		}

		if !f.isAdjust {
			t.adjustTimerAt(handle, f)
		}

		f.proc()
		f.lastCall = curTime
	}
}

func (t *Timer) doReg(callback func(), intervalSec time.Duration, oneTime bool) int {
	if callback == nil || intervalSec <= 0 {
		return 0
	}

	handle := t.handleSeed
	t.callbacks[handle] = &funcInfo{
		suspended: false,
		oneTime:   oneTime,
		isAdjust:  true,
		interval:  intervalSec,
		lastCall:  time.Now(),
		proc:      callback,
	}

	t.handleSeed++
	return handle
}

// RegTimerAt 添加一个整点定时任务 hour,min,sec:时分秒 oneTime:是否只执行一次
func (t *Timer) RegTimerAt(callback func(), hour, min, sec int, oneTime bool) int {
	return t.doRegTimerAt(callback, hour, min, sec, oneTime)
}

func (t *Timer) adjustTimerAt(handle int, f *funcInfo) {
	if !f.oneTime {
		f.interval = 86400 * time.Second
		f.isAdjust = true
	}
}

func (t *Timer) doRegTimerAt(callback func(), hour, min, sec int, oneTime bool) int {

	var (
		interval   time.Duration
		regTimeSec int64
		nowTimeSec int64
		localHour  int
		localMin   int
		localSec   int
	)
	localHour, localMin, localSec = time.Now().Local().Clock()
	regTimeSec = int64(3600*hour + 60*min + sec)                    // 注册的时间点距离今日0点过去的秒数
	nowTimeSec = int64(3600*localHour+60*localMin+localSec) % 86400 // 今天过去的秒数
	//注册的是今日时间ls
	if regTimeSec > nowTimeSec {
		interval = time.Duration(regTimeSec-nowTimeSec) * time.Second
	} else {
		interval = 86400*time.Second - time.Duration(nowTimeSec-regTimeSec)*time.Second
	}

	if callback == nil || interval <= 0 {
		return 0
	}

	handle := t.handleSeed
	t.callbacks[handle] = &funcInfo{
		suspended: false,
		oneTime:   oneTime,
		isAdjust:  false,
		interval:  interval,
		lastCall:  time.Now(),
		proc:      callback,
	}

	t.handleSeed++
	return handle
}
