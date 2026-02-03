package structure

import (
	. "github.com/codeshelldev/gotl/pkg/configutils/types"
)

type ENV struct {
	CONFIG_PATH   		string
	TOKENS_DIR    		string

	DEFAULTS_PATH 		string
	FAVICON_PATH  		string

	DB_PATH				string
	
	INSECURE      		bool

	TOKENS				[]string

	CONFIGS      		map[string]*CONFIG
}

type CONFIG struct {
	TYPE				ConfigType
	NAME				string						`koanf:"name"`
	SERVICE				SERVICE 					`koanf:"service"`
	API					API						    `koanf:"api"`
	// DEPRECATION overrides in Token Config
	SETTINGS      		SETTINGS					`koanf:"settings"        token>aliases:"overrides" token>onuse:".overrides>>deprecated"       deprecation:"{b,fg=yellow}\x60{s}overrides{/}\x60{/} is no longer needed in {b}Token Configs{/}\nUse {b,fg=green}\x60settings\x60{/} instead"`
}

type ConfigType string

const (
	TOKEN ConfigType = "token"
	MAIN ConfigType = "main"
)

type SERVICE struct {
	HOSTNAMES			Opt[[]string]				`koanf:"hostnames"       env>aliases:".hostnames"`
	PORT				string						`koanf:"port"            env>aliases:".port"`
	LOG_LEVEL			string						`koanf:"loglevel"        env>aliases:".loglevel"`
}

type API struct {
	URL					URL							`koanf:"url"             env>aliases:".apiurl"`
	// DEPRECATION token, tokens in Token Config
	// DEPRECATION api.token => api.tokens
	TOKENS				[]string					`koanf:"tokens"          env>aliases:".apitokens,.apitoken" aliases:"token" token>aliases:".tokens,.token" token>onuse:".tokens,.token,token>>deprecated" onuse:"token>>deprecated" deprecation:".tokens,.token>>{b,fg=yellow}\x60{s}tokens{/}\x60{/} and {b,fg=yellow}\x60{s}token{/}\x60{/} will not be at {b}root{/} anymore\nUse {b,fg=green}\x60api.tokens\x60{/} instead|token>>{b,fg=yellow}\x60{s}api.token{/}\x60{/} will be {u}removed{/} in favor of {b,fg=green}\x60api.tokens\x60{/}"`																					
	AUTH				AUTH						`koanf:"auth"`
}

type AUTH struct {
	METHODS				Opt[[]string]				`koanf:"methods"         env>aliases:".authmethods"`
	// DEPRECATION auth.token => auth.tokens
	TOKENS				[]Token						`koanf:"tokens"          aliases:"token" onuse:"token>>deprecated" deprecation:"{b,fg=yellow}\x60{s}api.auth.token{/}\x60{/} will be removed\nUse {b,fg=green}\x60api.auth.tokens\x60{/} instead"`
}

type Token struct {
	Set					[]string					`koanf:"set"`
	Methods				[]string					`koanf:"methods"`
}

type SETTINGS struct {
	ACCESS 				ACCESS 						`koanf:"access"`
	MESSAGE				MESSAGE						`koanf:"message"`
}

type MESSAGE struct {
	VARIABLES         	Opt[map[string]any]			`koanf:"variables"       childtransform:"upper"`
	FIELD_MAPPINGS      Opt[map[string][]FieldMapping]`koanf:"fieldmappings"   childtransform:"default"`
	TEMPLATE  			Opt[string]					`koanf:"template"`
}

type FieldMapping struct {
	Field 				string 						`koanf:"field"`
	Score 				int    						`koanf:"score"`
}

type ACCESS struct {
	ENDPOINTS			Opt[AllowBlockSlice]		`koanf:"endpoints"`
	FIELD_POLICIES		Opt[map[string]FieldPolicies]`koanf:"fieldpolicies"   childtransform:"default"`
	RATE_LIMITING		Opt[RateLimiting]			`koanf:"ratelimiting"`
	IP_FILTER			Opt[AllowBlockSlice]		`koanf:"ipfilter"`
	TRUSTED_IPS			Opt[[]IPOrNet]				`koanf:"trustedips"`
	TRUSTED_PROXIES		Opt[[]IPOrNet]				`koanf:"trustedproxies"`
}

type FieldPolicy struct {
	Value				any						    `koanf:"value"`
	Action				string						`koanf:"action"`
}

type RateLimiting struct {
	Limit				int							`koanf:"limit"`
	Period				TimeDuration				`koanf:"period"`
}