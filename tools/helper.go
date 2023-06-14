package tools

import (
	"fmt"
	"math/rand"
	"runtime"
	"runtime/debug"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func RunFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

func Catch() {
	if r := recover(); r != nil {
		Log.Error(fmt.Sprintf("%v", r)) //输出panic信息
		Log.Error(string(debug.Stack()[:]))
	}
}

// [a,b]
func GetRandInt(a, b int) int {
	if a > b {
		a, b = b, a
	}
	return a + rand.Intn(b-a+1)
}

func GetRandInt32(a, b int32) int32 {
	if a > b {
		a, b = b, a
	}
	return a + rand.Int31n(b-a+1)
}

var randStringPool = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "N", "M", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "n", "m", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	"1", "2", "3", "4", "6", "7", "8", "9"}

func GetRandString() string {
	poolLen := len(randStringPool)
	randomString := ""
	for i := 0; i != 6; i++ {
		randomString = randomString + randStringPool[GetRandInt(0, poolLen-1)]
	}
	return randomString
}

func GetRandPass(length int) string {
	if length <= 0 {
		Log.Panic("Password length error")
	}
	baseStr := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~!@#$%^&*()_+[{]};:',<.>/?"
	r := rand.New(rand.NewSource(time.Now().UnixNano() + rand.Int63()))
	bytes := make([]byte, length)
	l := len(baseStr)
	for i := 0; i < length; i++ {
		bytes[i] = baseStr[r.Intn(l)]
	}
	return string(bytes)
}

func MergeMaps(map1, map2 map[string]interface{}) map[string]interface{} {
	// 创建一个新的 map 用于存放结果
	result := make(map[string]interface{})

	// 将 map1 的内容复制到结果中
	for k, v := range map1 {
		result[k] = v
	}

	// 将 map2 的内容复制到结果中，如果有重复的键，那么会覆盖 map1 中的值
	for k, v := range map2 {
		result[k] = v
	}

	return result
}
