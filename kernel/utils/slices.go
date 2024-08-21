package utils

import (
	"fmt"
	"strconv"
	"strings"
)

type Float64Slice []float64
type IntSlice []int
type StringSlice []string
type IntMap map[int]int

func IntSliceFromString(str string, sep string) (IntSlice, error) {
	if len(str) == 0 {
		return IntSlice(make([]int, 0)), nil
	}
	strs := strings.Split(str, sep)
	var err error
	var result = make(IntSlice, len(strs))
	for i := 0; i < len(strs); i++ {
		if len(strs[i]) == 0 {
			continue
		}
		result[i], err = strconv.Atoi(strs[i])
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func Float64SliceFromString(str string, sep string) (Float64Slice, error) {
	if len(str) == 0 {
		return Float64Slice(make([]float64, 0)), nil
	}
	strs := strings.Split(str, sep)
	var err error
	var result = make(Float64Slice, len(strs))
	for i := 0; i < len(strs); i++ {
		if len(strs[i]) == 0 {
			continue
		}
		result[i], err = strconv.ParseFloat(strs[i], 64)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (this IntSlice) Index(element int) int {
	for i, v := range this {
		if v == element {
			return i
		}
	}
	return -1
}

func (this IntSlice) RemoveIndex(index int) IntSlice {
	if index < 0 || index >= len(this) {
		return this
	}
	return append(this[:index], this[index+1:]...)
}

func (this IntSlice) RemoveElement(element int) IntSlice {
	for i, v := range this {
		if v == element {
			return append(this[:i], this[i+1:]...)
		}
	}
	return this
}

func (this IntSlice) Add(element int) IntSlice {
	return append(this, element)
}

func (this IntSlice) AddUnique(element int) IntSlice {
	if this.Index(element) < 0 {
		return this
	}
	return append(this, element)
}

func (this IntSlice) String(sep string) string {
	var arrStr = make([]string, len(this))
	for i, v := range this {
		arrStr[i] = strconv.Itoa(v)
	}
	return strings.Join(arrStr, sep)
}

func (this IntSlice) Len() int {
	return len(this)
}

// 玩家按战力排名
func (this IntSlice) Less(i, j int) bool {
	if this[j] != this[i] {
		return this[j] > this[i]
	}
	return false
}

func (this IntSlice) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func ConvertInt32SliceToIntSlice(origin []int32) []int {
	ret := make([]int, len(origin))
	for i, v := range origin {
		ret[i] = int(v)
	}
	return ret
}

func ConvertIntSlice2Int32Slice(origin []int) []int32 {
	ret := make([]int32, len(origin))
	for i, v := range origin {
		ret[i] = int32(v)
	}
	return ret
}

func ConvertMapIntToInt32(origin map[int]int) map[int32]int32 {
	ret := make(map[int32]int32, len(origin))
	for i, v := range origin {
		ret[int32(i)] = int32(v)
	}
	return ret
}

func ConvertMapInt32ToInt(origin map[int32]int32) map[int]int {
	ret := make(map[int]int, len(origin))
	for i, v := range origin {
		ret[int(i)] = int(v)
	}
	return ret
}

func SliceIntUnique(origin []int) []int {
	ret := make([]int, 0)
	tempMap := make(map[int]struct{})
	for _, v := range origin {
		if _, ok := tempMap[v]; ok {
			continue
		}
		ret = append(ret, v)
		tempMap[v] = struct{}{}
	}
	return ret
}
func SliceInt32Unique(origin []int32) []int32 {
	ret := make([]int32, 0)
	tempMap := make(map[int32]struct{})
	for _, v := range origin {
		if _, ok := tempMap[v]; ok {
			continue
		}
		ret = append(ret, v)
		tempMap[v] = struct{}{}
	}
	return ret
}

// 2维int转string
func SliceInt2ToString(arr [][]int, sep1 string, sep2 string) string {

	slice1 := make([]string, len(arr))
	for k, v := range arr {
		slice1[k] = JoinIntSlice(v, sep1)
	}
	return strings.Join(slice1, sep2)
}

// 2维int转1维string
func SliceInt2ToSliceString1(arr [][]int, sep1 string) []string {

	slice1 := make([]string, len(arr))
	for k, v := range arr {
		slice1[k] = JoinIntSlice(v, sep1)
	}
	return slice1
}

func JoinIntSlice(a []int, sep string) string {
	l := len(a)
	if l == 0 {
		return ""
	}
	b := make([]string, l)
	for i, v := range a {
		b[i] = strconv.Itoa(v)
	}
	return strings.Join(b, sep)
}

func JoinInt32Slice(a []int32, sep string) string {
	l := len(a)
	if l == 0 {
		return ""
	}
	b := make([]string, l)
	for i, v := range a {
		b[i] = strconv.Itoa(int(v))
	}
	return strings.Join(b, sep)
}

func InterfaceSlice2StringSlice(arr []interface{}) []string {
	strArr := make([]string, len(arr))
	for index, v := range arr {
		strArr[index] = fmt.Sprint(v)
	}
	return strArr
}

// intmap -> string
func IntMap2ToString(arr IntMap) string {
	slice1 := make([]string, 0)
	for k, v := range arr {
		slice1 = append(slice1, fmt.Sprintf("%d,%d", k, v))
	}
	return strings.Join(slice1, ";")
}
