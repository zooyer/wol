package wol

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
)

func ListenWOL(port int) error {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: port,
	})
	if err != nil {
		return err
	}
	defer conn.Close()

	var buf = make([]byte, 108)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return err
		}
		fmt.Println("---------- wol server receive wol begin ---------")
		for i := 0; i < n; i += 6 {
			for j := 0; j < 6; j++ {
				fmt.Printf("%#02X ", buf[i+j])
			}
			fmt.Println()
		}
		fmt.Println("---------- wol server receive wol end -----------")
	}
}

// WOL wol唤醒mac主机
func WOL(mac string, port ...int) (err error) {
	packet, err := MagicPacket(mac)
	if err != nil {
		return
	}
	port = append(port, 7, 9)
	list := broadcastAddresses()
	if len(list) == 0 {
		return errors.New("not available interface")
	}

	var addr, laddr net.UDPAddr
	laddr.IP = net.IPv4zero

	for _, p := range port {
		for _, ip := range list {
			addr.IP = net.ParseIP(ip)
			addr.Port = p
			var conn *net.UDPConn
			if conn, err = net.DialUDP("udp", &laddr, &addr); err != nil {
				return
			}
			if _, err = io.Copy(conn, bytes.NewBuffer(packet)); err != nil {
				return
			}
			if err = conn.Close(); err != nil {
				return
			}
		}
	}

	return
}

// MagicPacket 封装魔术包
func MagicPacket(mac string) ([]byte, error) {
	mac = strings.ReplaceAll(mac, ":", "")
	mac = strings.ReplaceAll(mac, "-", "")
	mac = strings.ToUpper(mac)

	var head = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	var tail = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	body, err := hex.DecodeString(mac)
	if err != nil {
		return nil, err
	}

	if len(body) != 6 {
		return nil, errors.New("mac address wrong")
	}

	var packet = make([]byte, 0, 108)

	packet = append(packet, head...)
	for i := 0; i < 16; i++ {
		packet = append(packet, body...)
	}
	packet = append(packet, tail...)

	return packet, nil
}

func ips() []string {
	var list []string
	address, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	for _, addr := range address {
		if ip, ok := addr.(*net.IPNet); ok && ip.IP.IsGlobalUnicast() && ip.IP.To4() != nil {
			list = append(list, ip.IP.To4().String())
		}
	}

	return list
}

func broadcastAddresses() (list []string) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}

	for _, iface := range interfaces {
		// 跳过没有启用的接口
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		// 获取接口的地址信息
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		// 遍历每个接口的地址
		for _, addr := range addrs {
			if ip, ok := addr.(*net.IPNet); ok && ip.IP.IsGlobalUnicast() && ip.IP.To4() != nil {
				list = append(list, getBroadcastAddress(ip.IP, ip.Mask).To4().String())
			}
		}
	}

	return
}

// getBroadcastAddress 计算广播地址
func getBroadcastAddress(ip net.IP, mask net.IPMask) net.IP {
	// 网络地址: ip & mask
	network := ip.Mask(mask)

	// 广播地址: 将主机部分设置为 1
	for i := range network {
		network[i] |= ^mask[i]
	}

	return network
}
