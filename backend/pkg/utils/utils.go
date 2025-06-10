package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

func NewActionTimestamp() int64 {
	return time.Now().UnixMilli()
}

func PrintJson(data any) {
	value, _ := json.MarshalIndent(data, "", "    ")
	fmt.Printf("%s\n", value)
}

func NewGroupId() string {
	return GenerateRandomStringWithLength(8)
}

func NewServiceId(groupId string) string {
	return fmt.Sprintf("%s%s", groupId, GenerateRandomStringWithLength(8))
}

func NewVersionId() string {
	return GenerateRandomStringWithLength(5)
}

func HostIP() []net.IP {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Failed to get network interfaces:", err)
		return nil
	}

	var ipAddresses []net.IP
	for _, iface := range interfaces {
		// 跳过回环接口和未启用的接口
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		// 跳过 Docker 虚拟网卡
		if DockerInterface(iface.Name) {
			continue
		}

		addresses, err := iface.Addrs()
		if err != nil {
			fmt.Printf("Failed to get addresses for interface %s: %v\n", iface.Name, err)
			continue
		}

		for _, addr := range addresses {
			// 检查是否为 IPv4 地址
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					ipAddresses = append(ipAddresses, ipNet.IP)
				}
			}
		}
	}

	return ipAddresses
}

// 判断是否为 Docker 虚拟网卡
func DockerInterface(ifaceName string) bool {
	// Docker 虚拟网卡通常以以下前缀开头
	dockerPrefixes := []string{"docker", "veth", "br-", "cni", "flannel"}
	for _, prefix := range dockerPrefixes {
		if len(ifaceName) >= len(prefix) && ifaceName[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}
