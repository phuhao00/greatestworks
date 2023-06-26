package idgenerator

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

func FetchServiceID(svcName string, svcAddr string) string {
	host, port, err := net.SplitHostPort(svcAddr)
	if err != nil {
		return svcName + "-" + svcAddr
	}
	ipInt := ip2Int(host)
	return fmt.Sprintf("%v-%v", ipInt, port)
}

func ip2Int(ip string) uint32 {
	var ipInt uint32
	err := binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &ipInt)
	if err != nil {
		return 0
	}
	return ipInt
}
