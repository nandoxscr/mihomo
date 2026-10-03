//go:build xclashrev

package executor

import (
	"runtime"

	"github.com/metacubex/mihomo/config"
	"github.com/metacubex/mihomo/log"
	"github.com/metacubex/mihomo/tunnel"
)

// ApplyConfigOffline initializes proxies, providers, and DNS without starting
// inbound listeners (HTTP, SOCKS, Mixed), TUN devices, or IPTables rules.
// This allows querying proxies, groups, providers, and testing latency
// before the VPN/proxy tunnel service is started.
func ApplyConfigOffline(cfg *config.Config) {
	mux.Lock()
	defer mux.Unlock()
	log.SetLevel(cfg.General.LogLevel)

	tunnel.OnSuspend()

	updateExperimental(cfg.Experimental)
	updateUsers(cfg.Users)
	updateProxies(cfg.Proxies, cfg.Providers)
	updateGeneral(cfg.General, false)
	updateDNS(cfg.DNS, cfg.General.IPv6)

	tunnel.OnInnerLoading()

	initInnerTcp()
	loadProvider(cfg.Providers)
	updateProfile(cfg)
	runtime.GC()
	tunnel.OnRunning()
}
