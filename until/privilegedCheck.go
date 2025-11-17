package until

import (
	"log"
	"os"
	"strings"
)

// 检查 CapEff 是否为 ffffffffffffffff
func hasAllCaps() bool {
	data, err := os.ReadFile("/proc/self/status")
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "CapEff:") {
			val := strings.TrimSpace(strings.TrimPrefix(line, "CapEff:"))
			return strings.HasPrefix(val, "ffffffff") // full caps
		}
	}
	return false
}

// 检查是否可写入特权文件
func canWritePrivilegedFiles() bool {
	// /dev/kmsg（普通容器不可写）
	if f, err := os.OpenFile("/dev/kmsg", os.O_WRONLY, 0o0); err == nil {
		f.Close()
		return true
	}

	// /proc/sysrq-trigger（普通容器不可写）
	if f, err := os.OpenFile("/proc/sysrq-trigger", os.O_WRONLY, 0o0); err == nil {
		f.Close()
		return true
	}
	return false
}

// 检查是否能读取受限内核参数
func canAccessRestrictedKernelInfo() bool {
	_, err := os.ReadFile("/proc/kcore") // 普通容器不可读
	return err == nil
}

// ===============================
// 最终接口：一行调用即可
// ===============================
func IsPrivileged() bool {
	log.Println("运行环境检查...")
	// 一层：capabilities
	if hasAllCaps() {
		return true
	}
	// 二层：写敏感设备
	if canWritePrivilegedFiles() {
		return true
	}
	// 三层：读取内核核心区
	if canAccessRestrictedKernelInfo() {
		return true
	}
	return false
}
