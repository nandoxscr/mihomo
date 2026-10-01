//go:build with_ebpf && android && xclashrev

package sing_ebpf

import (
	"github.com/metacubex/mihomo/listener/sing_tun"
	tun "github.com/metacubex/sing-tun"
)

func androidPackageManager() (tun.PackageManager, error) {
	return sing_tun.GetPackageManager()
}
