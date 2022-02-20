package help

func TakeWhile(str string, del rune) string {
	var result string
	for _, r := range str {
		if r == del {
			break
		}
		result += string(r)
	}
	return result
}

func Reverse(str string) string {
	var result string
	for _, r := range str {
		result = string(r) + result
	}
	return result
}
