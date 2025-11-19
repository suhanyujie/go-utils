package times

import "time"

func GetNowMillisecond() int64 {
	return time.Now().UnixNano() / 1e6
}

func GetMillisecond(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

func GetNowSecond() int64 {
	return time.Now().Unix()
}

func GetNowNanoSecond() int64 {
	return time.Now().UnixNano()
}

func GetBeiJingTime() time.Time {
	timelocal, _ := time.LoadLocation("Asia/Chongqing")
	time.Local = timelocal
	return time.Now().Local()
}
