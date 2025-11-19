package randx

import (
	"math/rand"
	"time"

	"github.com/shopspring/decimal"
)

func GenIntWithRange(min, max int) int {
	rand.Seed(time.Now().Unix()) //随机种子
	return rand.Intn(max-min) + min
}

// 针对一个小数位的浮点数生成，eg: [1.1, 1.5]
// [min, max]
func GenFloatWithRange[F float32 | float64](min, max F) F {
	rand.Seed(time.Now().UnixNano()) //随机种子
	precisionNum := 10
	minInt := int(min * F(precisionNum))
	maxInt := int(max*F(precisionNum)) + 1 // for range [min, max]
	resInt := GenIntWithRange(minInt, maxInt)
	resultF := decimal.NewFromFloat(float64(resInt)).DivRound(decimal.NewFromInt(int64(precisionNum)), 2).InexactFloat64()
	return F(resultF)
}

// 生成 count 个 [start,end) 结束的不重复的随机数
// ref: https://blog.csdn.net/books1958/article/details/44923779
func GenRandomNumberByCnt(start int, end int, count int) []int {
	//范围检查
	if end < start || (end-start) < count {
		return nil
	}
	//存放结果的slice
	nums := make([]int, 0)
	//随机数生成器，加入时间戳保证每次生成的随机数不一样
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(nums) < count {
		//生成随机数
		num := r.Intn((end - start)) + start
		//查重
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}
		if !exist {
			nums = append(nums, num)
		}
	}

	return nums
}
