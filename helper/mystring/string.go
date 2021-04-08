package mystring

// 元素为 string 的 slice 版本的 InArray 函数
func StringInArray(needle string, strSlice []string) bool {
	hasFound := false
	if len(strSlice) < 1 {
		return hasFound
	}
	for _, val := range strSlice {
		if needle == val {
			hasFound = true
			break
		}
	}
	return hasFound
}
