package calc

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func calc() {
	NowTime := time.Now()
	src := "5:00:00"
	timeStrings := strings.Split(src, ":")
	if len(timeStrings) != 3 {
		return
	}
	hour, err := strconv.Atoi(timeStrings[0])
	if err != nil {
		return
	}
	minute, err := strconv.Atoi(timeStrings[1])
	if err != nil {
		return
	}
	second, err := strconv.Atoi(timeStrings[2])
	if err != nil {
		return
	}
	fillTime := int64((int32(hour) * 3600) + (int32(minute) * 60) + int32(second))
	fmt.Println(fillTime)
	Today := time.Unix(NowTime.Unix()-fillTime, 0)
	fmt.Println(Today.Day())
	t := 1 << (uint32(Today.Day() - 1))
	fmt.Println("int32类型位移后:", t)
	fmt.Println("一共有:", NumberOf1InBinary(t), "个1")
}

// 2进制数字有多少个1
func NumberOf1InBinary(n int) int {
	count := 0
	for n != 0 {
		count++
		n = n & (n - 1)
	}
	return count
}
