package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// MD5 md5加密
func MD5(src string) string {
	w := md5.New()
	w.Write([]byte(src))
	return hex.EncodeToString(w.Sum(nil))
}

// GUID 产生新的GUID
func GUID() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	strMD5 := base64.URLEncoding.EncodeToString(b)
	return MD5(strMD5)
}

// Token 产生新的用户登录验证码
func Token() string {
	return fmt.Sprint(time.Now().Unix(), ":", GUID())
}

func Int32ArrayToString(src []int32, flag string) (out string) {
	// 没有的就直接返回
	if len(src) == 0 {
		return ""
	}
	out = ""
	for k, v := range src {
		if k == len(src)-1 {
			out = fmt.Sprint(out, v)
		} else {
			out = fmt.Sprint(out, v, flag)
		}
	}

	return
}

func StringToInt32Array(src, flag string) (out []int32) {
	// 没有的就直接返回
	if src == "" {
		return nil
	}

	strs := strings.Split(src, flag)
	for _, v := range strs {
		data, err := strconv.Atoi(v)
		if err != nil {
			return nil
		}

		out = append(out, int32(data))
	}

	return
}

func StringArrayToInt32Array(src []string) (out []int32) {
	for _, v := range src {
		data, err := strconv.Atoi(v)
		if err != nil {
			return nil
		}
		out = append(out, int32(data))
	}

	return
}

func StrToInt32(src string) int32 {
	data, err := strconv.Atoi(src)
	if err != nil {
		return 0
	}

	return int32(data)
}

func StrToFloat(src string) float64 {
	data, err := strconv.ParseFloat(src, 32)
	if err != nil {
		return 0
	}

	return data
}

func StrToInt64(src string) int64 {
	data, err := strconv.ParseInt(src, 10, 64)
	if err != nil {
		return 0
	}

	return data
}

func StrToUInt64(src string) uint64 {
	data, err := strconv.ParseInt(src, 10, 64)
	if err != nil {
		return 0
	}

	return uint64(data)
}

func StrToInt(src string) int {
	data, err := strconv.Atoi(src)
	if err != nil {
		return 0
	}

	return data
}

// 表情解码
func UnicodeEmojiDecode(s string) string {
	//emoji表情的数据表达式
	re := regexp.MustCompile("\\[[\\\\u0-9a-zA-Z]+\\]")
	//提取emoji数据表达式
	reg := regexp.MustCompile("\\[\\\\u|]")
	src := re.FindAllString(s, -1)
	for i := 0; i < len(src); i++ {
		e := reg.ReplaceAllString(src[i], "")
		p, err := strconv.ParseInt(e, 16, 32)
		if err == nil {
			s = strings.Replace(s, src[i], string(rune(p)), -1)
		}
	}
	return s
}

// 表情转换
func UnicodeEmojiCode(s string) string {
	ret := ""
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if len(string(rs[i])) == 4 {
			u := `[\u` + strconv.FormatInt(int64(rs[i]), 16) + `]`
			ret += u

		} else {
			ret += string(rs[i])
		}
	}
	return ret
}

// 删除空切片(字节)
func TrimSpace(s []byte) []byte {
	b := s[:0]
	for _, x := range s {
		if x != ' ' {
			b = append(b, x)
		}
	}
	return b
}

/**
 * 通过string获得一个[]string
 * @param str  10,1000,1000,1000
 * @param sep1 分隔符 ","
 */
func NewStringSlice(str string, sep string) []string {
	intSliceList := make([]string, 0)
	list := strings.Split(str, sep)
	for _, v := range list {
		intSliceList = append(intSliceList, v)
	}
	return intSliceList
}
