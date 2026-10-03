//go:build xclashrev

package hub

import (
	"github.com/metacubex/mihomo/config"
	"github.com/metacubex/mihomo/hub/executor"
)

// ApplyConfigOffline dispatches configuration only to proxies, providers,
// and DNS without starting inbound listeners or TUN.
func ApplyConfigOffline(cfg *config.Config) {
	executor.ApplyConfigOffline(cfg)
}

// ParseOffline parses config and applies it offline (without listeners or TUN).
func ParseOffline(configBytes []byte, options ...Option) error {
	var cfg *config.Config
	var err error

	if len(configBytes) != 0 {
		cfg, err = executor.ParseWithBytes(configBytes)
	} else {
		cfg, err = executor.Parse()
	}

	if err != nil {
		return err
	}

	for _, option := range options {
		option(cfg)
	}

	ApplyConfigOffline(cfg)
	return nil
}
