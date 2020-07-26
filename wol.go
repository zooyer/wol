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
func WOL(mac string, port ...int) error {
	packet, err := MagicPacket(mac)
	if err != nil {
		return err
	}
	port = append(port, 7, 9)
	list := ips()
	if len(list) == 0 {
		return errors.New("not available interface")
	}

	var addr, laddr net.UDPAddr
	addr.IP = net.ParseIP("255.255.255.255")

	for _, p := range port {
		addr.Port = p
		for _, ip := range list {
			laddr.IP = net.ParseIP(ip)
			wol, err := net.DialUDP("udp", &laddr, &addr)
			if err != nil {
				return err
			}
			if _, err = io.Copy(wol, bytes.NewBuffer(packet)); err != nil {
				return err
			}
			if err = wol.Close(); err != nil {
				return err
			}
		}
	}

	return nil
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
