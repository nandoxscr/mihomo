//go:build xclashrev

package executor

import (
	"runtime"

	"github.com/metacubex/mihomo/config"
	"github.com/metacubex/mihomo/log"
	"github.com/metacubex/mihomo/tunnel"
)

// ApplyConfigOffline initializes proxies, providers, and profile without starting
// inbound listeners (HTTP, SOCKS, Mixed), TUN devices, or IPTables rules.
func ApplyConfigOffline(cfg *config.Config) {
	mux.Lock()
	defer mux.Unlock()
	if cfg.General != nil {
		log.SetLevel(cfg.General.LogLevel)
	}

	tunnel.OnSuspend()

	if cfg.Experimental != nil {
		updateExperimental(cfg.Experimental)
	}
	updateProxies(cfg.Proxies, cfg.Providers)
	updateRules(nil, nil, cfg.RuleProviders)
	if cfg.General != nil {
		updateGeneral(cfg.General, false)
	}
	if cfg.DNS != nil {
		updateDNS(cfg.DNS, cfg.General != nil && cfg.General.IPv6)
	}

	tunnel.OnInnerLoading()

	initInnerTcp()
	if cfg.Profile != nil {
		updateProfile(cfg)
	}
	go loadProvider(cfg.Providers)
	go loadProvider(cfg.RuleProviders)
	runtime.GC()
	tunnel.OnRunning()
}
