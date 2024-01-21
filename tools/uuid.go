package tools

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/sony/sonyflake"
	"net"
	"sync"
	"time"
)

var (
	uuidUint64Pool sync.Pool
	settings       sonyflake.Settings
)

func init() {
	f := func() (uint16, error) {
		var id uint16
		ifts, err := net.Interfaces()
		if err != nil {
			return 0, err
		}
		for _, ift := range ifts {
			if len(ift.HardwareAddr) >= 6 {
				id = uint16(ift.HardwareAddr[4])<<8 + uint16(ift.HardwareAddr[5])
				return id, nil
			}
		}
		return 0, fmt.Errorf("no suitable interface found")
	}
	settings = sonyflake.Settings{StartTime: time.Now(), MachineID: f}
	uuidUint64Pool = sync.Pool{
		New: func() interface{} {
			return sonyflake.NewSonyflake(settings)
		},
	}
	uuid.EnableRandPool()
}

func UuidUint64() uint64 {
	flake := uuidUint64Pool.Get().(*sonyflake.Sonyflake)
	if uid, err := flake.NextID(); err != nil {
		Log.Error(err.Error())
		return 0
	} else {
		return uid
	}
}

//UUID v1（基于时间和节点 ID）
//生成过程：UUID v1 使用当前时间（精确到100纳秒）、时钟序列和节点ID（通常是机器的MAC地址）来生成 UUID。时间和时钟序列保证了时间上的唯一性，而节点ID保证了空间上的唯一性。
//全球唯一性：理论上，由于时间的单调递增性和节点ID的唯一性，UUID v1 能够保证全球唯一。但它可能会暴露时间信息和节点（如服务器）的MAC地址，可能不适用于需要隐私或安全性的场景。
//UUID v4（随机生成）
//生成过程：UUID v4 完全基于随机数或伪随机数生成。它由 122 位随机数字和 6 位固定位组成（用于指示 UUID 版本和变体）。
//全球唯一性：虽然理论上 UUID v4 有极小的重复几率，但由于其高度随机性，实际应用中几乎可以保证全球唯一。UUID v4 是最常用的版本，因为它既简单又高效，不暴露任何系统或时间信息。
//UUID v3 和 v5（基于哈希和命名空间）
//生成过程：这两个版本的 UUID 是通过散列一组给定的命名空间和名称来生成的。UUID v3 使用 MD5 哈希，而 UUID v5 使用更安全的 SHA-1 哈希。
//全球唯一性：这些 UUID 的唯一性依赖于命名空间和名称的唯一性。如果输入的命名空间和名称组合是唯一的，那么生成的 UUID 也是唯一的。这些版本适用于需要根据固定名称生成固定 UUID 的场景。
//总结
//UUID v1 能够保证全球唯一性，但可能会暴露时间和节点信息。
//UUID v4 是最常用的，提供高度的随机性和全球唯一性，但不携带任何特定信息。
//UUID v3 和 v5 用于特定场景，它们的唯一性依赖于输入的命名空间和名称。

func UuidV1() string {
	u, err := uuid.NewUUID()
	if err != nil {
		Log.Error(err.Error())
		return ""
	}
	return u.String()
}

func UuidV3(data []byte) string {
	return uuid.NewMD5(uuid.NameSpaceDNS, data).String()
}

func UuidV5(data []byte) string {
	return uuid.NewSHA1(uuid.NameSpaceDNS, data).String()
}

func UuidV4() string {
	uuid.EnableRandPool()
	u, err := uuid.NewRandom()
	if err != nil {
		Log.Error(err.Error())
		return ""
	}
	return u.String()
}

func UuidV7() string {
	u, err := uuid.NewV7()
	if err != nil {
		Log.Error(err.Error())
		return ""
	}
	return u.String()
}
