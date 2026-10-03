//go:build xclashrev

package config

// ParseConfigOffline parses proxies, proxy-groups, proxy-providers,
// and rule-providers from raw config bytes for offline browsing and latency testing,
// skipping rules, DNS, and TUN initialization.
func ParseConfigOffline(buf []byte) (*Config, error) {
	rawCfg, err := UnmarshalRawConfig(buf)
	if err != nil {
		return nil, err
	}

	general, err := parseGeneral(rawCfg)
	if err != nil {
		return nil, err
	}

	proxies, providers, err := parseProxies(rawCfg)
	if err != nil {
		return nil, err
	}

	ruleProviders, err := parseRuleProviders(rawCfg)
	if err != nil {
		return nil, err
	}

	profile, _ := parseProfile(rawCfg)

	return &Config{
		General:       general,
		Profile:       profile,
		Proxies:       proxies,
		Providers:     providers,
		RuleProviders: ruleProviders,
	}, nil
}
