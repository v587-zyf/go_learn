package calc

import "fmt"

// 快速排序算法
func QuickSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	splitData := arr[0]          //第一个数据
	low := make([]int, 0, 0)     //比我小的数据
	height := make([]int, 0, 0)  //比我大的数据
	mid := make([]int, 0, 0)     //与我一样大的数据
	mid = append(mid, splitData) //加入一个
	for i := 1; i < len(arr); i++ {
		if arr[i] < splitData {
			low = append(low, arr[i])
		} else if arr[i] > splitData {
			height = append(height, arr[i])
		} else {
			mid = append(mid, arr[i])
		}
	}
	low, height = QuickSort(low), QuickSort(height)
	myArr := append(append(low, mid...), height...)
	return myArr
}

func TestQuickSort() {
	arr := []int{1, 9, 10, 30, 2, 5, 45, 8, 63, 234, 12}
	fmt.Println(QuickSort(arr))
}
