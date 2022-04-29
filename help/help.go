package help

import (
	"fmt"
	"os"
	"strconv"
)

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

func GetFormat(name string) string {
	return Reverse(TakeWhile(Reverse(name), '.'))
}

func GetSize(name string) (string, error) {
	file, err := os.Open(name)
	if err != nil {
		return "", err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return "", err
	}
	return normalize(float32(stat.Size())), nil
}

func normalize(memory float32) string {
	switch {
	case memory < 1024:
		return fmt.Sprintf("%.2f", memory) + "b"
	case memory < 1048576:
		return fmt.Sprintf("%.2f", memory/1024) + "kb"
	case memory < 1073741824:
		return fmt.Sprintf("%.2f", memory/1024/1024) + "mb"
	case memory < 1099511627776:
		return fmt.Sprintf("%.2f", memory/1024/1024/1024) + "gb"
	default:
		return strconv.Itoa(int(memory)) + "b"
	}
}
