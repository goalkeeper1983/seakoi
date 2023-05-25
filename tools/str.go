package tools

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//零GC 类型转换
//func Str2Bytes(s string) []byte {
//	x := (*[2]uintptr)(unsafe.Pointer(&s))
//	h := [3]uintptr{x[0], x[1], x[1]}
//	return *(*[]byte)(unsafe.Pointer(&h))
//}
//
//func Bytes2Str(b []byte) string {
//	return *(*string)(unsafe.Pointer(&b))
//}

func ToString(i interface{}) string {
	return fmt.Sprintf("%+v", i)
}

// TO json return string
func TOJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// TO json return byte
func TOJSONByte(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

// 首字幕大写
func Capitalize(str string) string {
	if len(str) == 0 {
		return ""
	}
	vv := []byte(str)
	if vv[0] >= 97 && vv[0] <= 122 { //ascii
		vv[0] -= 32
	}
	return string(vv)
}

func StringToUint(args string) uint {
	if args == "" {
		return 0
	}
	r, err := strconv.Atoi(args)
	if err != nil {
		return 0
	}
	return uint(r)
}

func Uint64ToString(args uint64) string {
	return strconv.FormatUint(args, 10)
}

func StringToUint64(args string) (uint64, error) {
	return strconv.ParseUint(args, 10, 64)
}

func Int64ToString(args int64) string {
	return strconv.FormatInt(args, 10)
}

func StringToInt(args string) int {
	if args == "" {
		return 0
	}
	r, err := strconv.Atoi(args)
	if err != nil {
		return 0
	}
	return r
}

func StringToInt32(args string) int32 {
	return int32(StringToInt(args))
}

func StringToInt64(args string) (int64, error) {
	return strconv.ParseInt(args, 10, 64)
}

func Float64ToString(f float64) string {
	return fmt.Sprintf("%v", f)
}

func Float32ToString(f float32) string {
	return fmt.Sprintf("%v", f)
}

// StringToFloat64 string到float64
func StringToFloat64(str string) float64 {
	if f, err := strconv.ParseFloat(str, 64); err == nil {
		return f
	} else {
		Log.Error(fmt.Sprintf("StringToFloat64Err:%v,%v", err.Error(), str))
		return 0
	}
}

func Int32ToString(i int32) string {
	return fmt.Sprintf("%v", i)
}

//StringToFloat32 string到float32
//func StringToFloat32(str string) float32 {
//	if f, err := strconv.ParseFloat(str, 32); err == nil {
//		return f
//	} else {
//		Log.Error(fmt.Sprintf("StringToFloat64Err:%v,%v", err.Error(), str))
//		return 0
//	}
//}

// 密码校验规则: 必须包含数字、大写字母、小写字母、特殊字符(如.@$!%*#_~?&^)至少3种的组合且长度在8-16之间
func VerifyLoginPassword(minLength, maxLength int, pwd string) bool {
	if len(pwd) < minLength || len(pwd) > maxLength {
		return false
	}
	// 过滤掉这四类字符以外的密码串,直接判断不合法
	re, err := regexp.Compile(`^[a-zA-Z0-9.@$!%*#_~?&^]{8,16}$`)
	if err != nil {
		return false
	}
	match := re.MatchString(pwd)
	if !match {
		return false
	}

	var level = 0
	//patternList := []string{`[0-9]+`, `[a-z]+`, `[A-Z]+`, `[.@$!%*#_~?&^]+`}//数字、大写字母、小写字母、特殊字符
	patternList := []string{`[0-9]+`, `[a-z]+`, `[A-Z]+`}
	for _, pattern := range patternList {
		match, _ = regexp.MatchString(pattern, pwd)
		if match {
			level++
		}
	}
	return level >= 3
}

func VerifySecondaryPassword(length int, pwd string) bool {
	if len(pwd) != length {
		return false
	}
	return IsDigit(pwd)
}

func IsDigit(s string) bool {
	pattern := "^[0-9]*$"
	match, _ := regexp.MatchString(pattern, s)
	return match
}

func handlingSpaces(input string) string {
	htmlSpecialCharsList := []string{"&", "\"", "'", "<", ">", "\\", "\t", "\n", "\r", " "}
	for _, v := range htmlSpecialCharsList {
		input = strings.Replace(input, v, "", -1)
	}
	return input
}
