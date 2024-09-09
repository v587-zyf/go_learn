package test

import (
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) {
	type testCase struct {
		str  string
		sep  string
		want []string
	}

	testGroup := map[string]testCase{
		"|": {str: "a|b|c|d", sep: "|", want: []string{"a", "b", "c", "d"}},
		",": {str: "a,b,c,d", sep: ",", want: []string{"a", "b", "c", "d"}},
		":": {str: "a:b:c:d", sep: ":", want: []string{"a", "b", "c", "d"}},
	}
	for name, testNode := range testGroup {
		t.Run(name, func(t *testing.T) {
			ret := Split(testNode.str, testNode.sep)
			if !reflect.DeepEqual(ret, testNode.want) {
				t.Errorf("got:%v, want:%v", ret, testNode.want)
			}
		})
	}
}

// go test -bench=Split
/**
 * goos: windows
 * goarch: amd64
 * pkg: demo/test
 * cpu: Intel(R) Core(TM) i5-6200U CPU @ 2.30GHz
 * BenchmarkSplit-4         4806472               258.0 ns/op
 * 4核 						4806472次  每次调用耗时258.0 ns(4806472次平均值)
 * PASS
 * ok      demo/test       1.717s
 */

// go test -bench=Split -benchmem
/**
 * goos: windows
 * goarch: amd64
 * pkg: demo/test
 * cpu: Intel(R) Core(TM) i5-6200U CPU @ 2.30GHz
 * BenchmarkSplit-4   4575782    			264.0 ns/op           		112 B/op          3 allocs/op
 * 4核              4575782次  每次调用耗时264.0 ns(4575782次平均值)	每次操作占用112字节	申请多少次内存
 * PASS
 * ok      demo/test       1.717s
 */
func BenchmarkSplit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Split("a,b,c,d", ",")
	}
}

func benchmarkFib(b *testing.B, n int) {
	for i := 0; i < b.N; i++ {
		Fib(n)
	}
}
func BenchmarkFib10(b *testing.B) { benchmarkFib(b, 10) }
func BenchmarkFib20(b *testing.B) { benchmarkFib(b, 20) }
