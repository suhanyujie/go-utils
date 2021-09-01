package slicex

import (
	"strconv"
	"strings"
)

// 查找在 arr1 中，但不在 arr2 中的元素
func ArrayDiff(arr1, arr2 []int64) (diffArr []int64) {
	if len(arr2) < 1 || len(arr1) < 1 {
		diffArr = arr1
		return
	}
	for i := 0; i < len(arr1); i++ {
		item := arr1[i]
		isIn := false
		for j := 0; j < len(arr2); j++ {
			if item == arr2[j] {
				isIn = true
				break
			}
		}
		if !isIn {
			diffArr = append(diffArr, item)
		}
	}
	return diffArr
}

// 去重，并且不会打乱顺序
func ArrayUnique(arr1 []int64) []int64 {
	resArr := []int64{}
	uniqueMap := map[int64]bool{}
	for _, one := range arr1 {
		if _, ok := uniqueMap[one]; ok {
			continue
		} else {
			uniqueMap[one] = true
			resArr = append(resArr, one)
		}
	}
	return resArr
}

func InArray(val int64, list []int64) bool {
	res := false
	for _, item := range list {
		if item == val {
			res = true
			break
		}
	}
	return res
}

// 两个 int64 数组切片的交集
func Int64Intersect(list1, list2 []int64) []int64 {
	uniqueMap := map[int64]bool{}
	for _, ele := range list1 {
		uniqueMap[ele] = true
	}
	result := make([]int64, 0)
	for _, ele := range list2 {
		if _, ok := uniqueMap[ele]; ok {
			result = append(result, ele)
		}
	}
	return result
}

// 移除 haystack 中的一些元素
func Int64RemoveSomeVal(haystack, removeVals []int64) []int64 {
	if len(removeVals) < 1 {
		return haystack
	}
	result := make([]int64, 0)
	for _, ele := range haystack {
		if InArray(ele, removeVals) {
			continue
		}
		result = append(result, ele)
	}
	return result
}

// Int64Explode 将切片内容拼接成字符串
func Int64Explode(list []int64, glue string) string {
	strList := make([]string, 0)
	for _, item := range list {
		strList = append(strList, strconv.FormatInt(item, 10))
	}
	return strings.Join(strList, glue)
}
