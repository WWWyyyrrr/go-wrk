package util

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

type HeaderList []string

func (i *HeaderList) String() string {
	out := []string{}
	for _, s := range *i {
		out = append(out, s)
	}
	return strings.Join(out, ", ")
}

func (i *HeaderList) Set(value string) error {
	*i = append(*i, value)
	return nil
}

// RedirectError specific error type that happens on redirection
type RedirectError struct {
	msg string
}

func (self *RedirectError) Error() string {
	return self.msg
}

func NewRedirectError(message string) *RedirectError {
	rt := RedirectError{msg: message}
	return &rt
}

// ByteSize a helper struct that implements the String() method and returns a human readable result. Very useful for %v formatting.
type ByteSize struct {
	Size float64
}

func (self ByteSize) String() string {
	var rt float64
	var suffix string
	const (
		Byte  = 1
		KByte = Byte * 1024
		MByte = KByte * 1024
		GByte = MByte * 1024
	)

	if self.Size > GByte {
		rt = self.Size / GByte
		suffix = "GB"
	} else if self.Size > MByte {
		rt = self.Size / MByte
		suffix = "MB"
	} else if self.Size > KByte {
		rt = self.Size / KByte
		suffix = "KB"
	} else {
		rt = self.Size
		suffix = "bytes"
	}

	srt := fmt.Sprintf("%.2f%v", rt, suffix)

	return srt
}

func MaxDuration(d1 time.Duration, d2 time.Duration) time.Duration {
	if d1 > d2 {
		return d1
	} else {
		return d2
	}
}

func MinDuration(d1 time.Duration, d2 time.Duration) time.Duration {
	if d1 < d2 {
		return d1
	} else {
		return d2
	}
}

// EstimateHttpHeadersSize had to create this because headers size was not counted
func EstimateHttpHeadersSize(headers http.Header) (result int64) {
	result = 0

	for k, v := range headers {
		result += int64(len(k) + len(": \r\n"))
		for _, s := range v {
			result += int64(len(s))
		}
	}

	result += int64(len("\r\n"))

	return result
}

func CheckIPAvailable(ip net.IP) (bool, error) {
	// 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		return false, err
	}
	// 遍历每个网络接口
	for _, iface := range interfaces {
		// 获取接口的所有IP地址
		addrs, err := iface.Addrs()
		if err != nil {
			return false, err
		}
		// 遍历每个IP地址并与给定的IP地址进行比较
		for _, addr := range addrs {
			var ipNet *net.IPNet
			var ok bool
			if ipNet, ok = addr.(*net.IPNet); !ok {
				continue
			}
			if ip.Equal(ipNet.IP) {
				return true, nil // 找到匹配的IP地址
			}
		}
	}
	return false, nil
}
