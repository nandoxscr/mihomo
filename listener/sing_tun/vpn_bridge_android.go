//go:build android && (!cmfa || xclashrev)

package sing_tun

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"net/netip"
	"sync"
	"time"

	"github.com/metacubex/mihomo/constant"
)

const vpnUidSocket = "@xyz.chz.xclashrev.vpn_uid"

type vpnUidClient struct {
	mu   sync.Mutex
	conn net.Conn
}

var globalVpnUidClient vpnUidClient

func (c *vpnUidClient) query(req []byte) (uint32, string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		conn, err := net.DialTimeout("unix", vpnUidSocket, 100*time.Millisecond)
		if err != nil {
			return 0, "", err
		}
		c.conn = conn
	}

	_ = c.conn.SetDeadline(time.Now().Add(250 * time.Millisecond))
	if _, err := c.conn.Write(req); err != nil {
		_ = c.conn.Close()
		c.conn = nil
		return 0, "", err
	}

	// Response format from XClashRevTunService VpnUidResolverServer:
	// 4 bytes: int32 uid (BigEndian)
	// 2 bytes: uint16 package name length (BigEndian)
	// N bytes: package name string
	var header [6]byte
	if _, err := io.ReadFull(c.conn, header[:]); err != nil {
		_ = c.conn.Close()
		c.conn = nil
		return 0, "", err
	}

	rawUid := int32(binary.BigEndian.Uint32(header[0:4]))
	pkgLen := int(binary.BigEndian.Uint16(header[4:6]))

	if rawUid <= 0 {
		return 0, "", errors.New("uid not found")
	}

	var pkgName string
	if pkgLen > 0 {
		pkgBuf := make([]byte, pkgLen)
		if _, err := io.ReadFull(c.conn, pkgBuf); err != nil {
			_ = c.conn.Close()
			c.conn = nil
			return uint32(rawUid), "", err
		}
		pkgName = string(pkgBuf)
	}

	return uint32(rawUid), pkgName, nil
}

func parseAddr(addr net.Addr) (netip.Addr, uint16, bool) {
	if addr == nil {
		return netip.Addr{}, 0, false
	}
	switch a := addr.(type) {
	case *net.TCPAddr:
		if ip, ok := netip.AddrFromSlice(a.IP); ok {
			return ip.Unmap(), uint16(a.Port), true
		}
	case *net.UDPAddr:
		if ip, ok := netip.AddrFromSlice(a.IP); ok {
			return ip.Unmap(), uint16(a.Port), true
		}
	}
	if ap, err := netip.ParseAddrPort(addr.String()); err == nil {
		return ap.Addr().Unmap(), ap.Port(), true
	}
	return netip.Addr{}, 0, false
}

func resolveFromVpnBridge(metadata *constant.Metadata) (uint32, string, error) {
	var protocol byte
	switch metadata.NetWork {
	case constant.TCP:
		protocol = 6
	case constant.UDP:
		protocol = 17
	default:
		return 0, "", errors.New("unsupported protocol")
	}

	var srcIP, dstIP netip.Addr
	var srcPort, dstPort uint16

	if ip, port, ok := parseAddr(metadata.RawSrcAddr); ok {
		srcIP, srcPort = ip, port
	} else if metadata.SrcIP.IsValid() {
		srcIP, srcPort = metadata.SrcIP.Unmap(), metadata.SrcPort
	}

	if ip, port, ok := parseAddr(metadata.RawDstAddr); ok {
		dstIP, dstPort = ip, port
	} else if metadata.DstIP.IsValid() {
		dstIP, dstPort = metadata.DstIP.Unmap(), metadata.DstPort
	}

	if !srcIP.IsValid() || !dstIP.IsValid() {
		return 0, "", errors.New("invalid IP addresses")
	}

	var req []byte
	if srcIP.Is4() && dstIP.Is4() {
		req = make([]byte, 14)
		req[0] = protocol
		req[1] = 4
		src4 := srcIP.As4()
		copy(req[2:6], src4[:])
		binary.BigEndian.PutUint16(req[6:8], srcPort)
		dst4 := dstIP.As4()
		copy(req[8:12], dst4[:])
		binary.BigEndian.PutUint16(req[12:14], dstPort)
	} else if srcIP.Is6() && dstIP.Is6() {
		req = make([]byte, 38)
		req[0] = protocol
		req[1] = 6
		src16 := srcIP.As16()
		copy(req[2:18], src16[:])
		binary.BigEndian.PutUint16(req[18:20], srcPort)
		dst16 := dstIP.As16()
		copy(req[20:36], dst16[:])
		binary.BigEndian.PutUint16(req[36:38], dstPort)
	} else {
		return 0, "", errors.New("mismatched IP families")
	}

	return globalVpnUidClient.query(req)
}
