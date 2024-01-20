package tools

import (
	"fmt"
	"math/rand"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

// 在多个协程（goroutine）中共享并使用同一个 rand.Rand 实例会带来几个潜在风险和问题，主要因为 rand.Rand 不是并发安全的。这些风险包括：
// 竞态条件（Race Condition）：
// 当多个协程同时读写同一个 rand.Rand 实例时，由于内部状态的并发修改，可能会导致竞态条件。这种情况下，随机数生成的结果可能会变得不可预测，甚至是错误的。
// 非确定性行为：
// 竞态条件可能导致生成的随机数序列每次运行时都不同，即使是在相同的初始条件下。这种非确定性可能会干扰程序的可测试性和可复现性。
// 性能问题：
// 为了避免竞态条件，你可能需要通过互斥锁（mutex）来同步对 rand.Rand 实例的访问。这种锁机制可能会引入性能瓶颈，特别是在高并发的环境中，由于多个协程频繁地争用同一个锁。
//
// 为了解决这些问题，有几种常用的方法：
// 为每个协程创建单独的 rand.Rand 实例：
// 这是最简单和最推荐的方法。每个协程使用独立的 rand.Rand 实例，这样就可以避免竞态条件和锁争用，同时保持每个协程的随机数生成独立和可预测。
// 使用互斥锁保护共享的 rand.Rand 实例：
// 如果你确实需要共享同一个 rand.Rand 实例，可以使用互斥锁来同步对它的访问。但是，这种方法可能会导致性能问题。
// 使用 sync.Pool 管理 rand.Rand 实例：
// sync.Pool 可以用来管理和重用 rand.Rand 实例，每个协程可以从池中获取一个实例，使用后再放回池中。这种方法可以在一定程度上减少创建新实例的开销，但仍然需要小心处理实例的同步问题。
var randPool sync.Pool

func init() {
	randPool = sync.Pool{
		New: func() interface{} {
			return rand.New(rand.NewSource(time.Now().UnixNano()))
		},
	}
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
	r := randPool.Get().(*rand.Rand) // 从池中获取一个rand.Rand实例
	defer randPool.Put(r)
	return a + r.Intn(b-a+1)
}

func GetRandInt32(a, b int32) int32 {
	if a > b {
		a, b = b, a
	}
	r := randPool.Get().(*rand.Rand) // 从池中获取一个rand.Rand实例
	defer randPool.Put(r)
	return a + r.Int31n(b-a+1)
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
	bytes := make([]byte, length)
	l := len(baseStr)
	for i := 0; i < length; i++ {
		bytes[i] = baseStr[GetRandInt32(0, int32(l))-1]
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
