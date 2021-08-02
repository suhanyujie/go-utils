package mystring

// StringInArray 元素为 string 的 slice 版本的 InArray 函数
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

// Substr 字符串截取
//start：正数 - 在字符串的指定位置开始,超出字符串长度强制把start变为字符串长度
//       负数 - 在从字符串结尾的指定位置开始
//       0 - 在字符串中的第一个字符处开始
//length:正数 - 从 start 参数所在的位置返回
//       负数 - 从字符串末端返回
func Substr(str string, start, length int) string {
	if length == 0 {
		return ""
	}
	runeStr := []rune(str)
	lenStr := len(runeStr)

	if start < 0 {
		start = lenStr + start
	}
	if start > lenStr {
		start = lenStr
	}
	end := start + length
	if end > lenStr {
		end = lenStr
	}
	if length < 0 {
		end = lenStr + length
	}
	if start > end {
		start, end = end, start
	}
	return string(runeStr[start:end])
}

