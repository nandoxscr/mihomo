//go:build !android || (cmfa && !xclashrev)

package sing_tun

import (
	"errors"

	"github.com/metacubex/mihomo/constant"
)

func resolveFromVpnBridge(metadata *constant.Metadata) (uint32, string, error) {
	return 0, "", errors.New("vpn bridge not supported on this platform")
}
