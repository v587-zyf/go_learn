package utils

import (
	"math"
	"math/rand"
	"time"
)

/**
*随机权重
*randData map[int]int{索引:权重，索引：权重}
*return 索引
 */
func RandWeightByMap(randData map[int]int) int {
	sum := 0
	for _, v := range randData {
		sum += v
	}
	if sum <= 0 {
		return -1
	}
	randNum := rand.Intn(sum)
	count := 0
	for k, v := range randData {
		count += v
		if randNum < count {
			return k
		}
	}
	return -1
}

func RandSeq(n int) string {
	var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// RoundFloat 使用 math.Round 将浮点数 f 四舍五入到小数点后 n 位。
func RoundFloat(f float64, n int) float64 {
	if f < 0 {
		return f
	}
	shift := math.Pow(10, float64(n))
	// 将浮点数乘以10的n次方，四舍五入到最近的整数，然后再除以10的n次方。
	return math.Round(f*shift) / shift
}

// 包含上下限 [min, max]
func RandomWithAll(min, max int) int64 {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return int64(rand.Intn(max-min+1) + min)
}

// 不包含上限 [min, max)
func RandomWithMin(min, max int) int64 {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return int64(rand.Intn(max-min) + min)
}

// 不包含下限 (min, max]
func RandomWithMax(min, max int) int64 {
	var res int64
	rand.New(rand.NewSource(time.Now().UnixNano()))
Restart:
	res = int64(rand.Intn(max-min+1) + min)
	if res == int64(min) {
		goto Restart
	}
	return res
}

// 都不包含 (min, max)
func RandomWithNo(min, max int) int64 {
	var res int64
	rand.New(rand.NewSource(time.Now().UnixNano()))
Restart:
	res = int64(rand.Intn(max-min) + min)
	if res == int64(min) {
		goto Restart
	}
	return res
}

// 向上取整float64
func CeilFloat64(x float64) int {
	return int(math.Ceil(RoundFloat(x, 2)))
}
