//go:build xclashrev

package hub

import (
	"github.com/metacubex/mihomo/config"
	"github.com/metacubex/mihomo/hub/executor"
)

// ApplyConfigOffline dispatches configuration only to proxies and providers
// without starting inbound listeners or TUN.
func ApplyConfigOffline(cfg *config.Config) {
	executor.ApplyConfigOffline(cfg)
}

// ParseOffline parses config and applies it offline (without listeners or TUN).
func ParseOffline(configBytes []byte) error {
	cfg, err := config.ParseConfigOffline(configBytes)
	if err != nil {
		return err
	}
	ApplyConfigOffline(cfg)
	return nil
}
