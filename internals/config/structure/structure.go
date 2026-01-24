package structure

type ENV struct {
	CONFIG_PATH   		string
	DEFAULTS_PATH 		string
	FAVICON_PATH  		string
	TOKENS_DIR    		string
	INSECURE      		bool

	CONFIGS      		map[string]*CONFIG
}

type CONFIG struct {
	NAME				string						`koanf:"name"`
	SERVICE				SERVICE 					`koanf:"service"`
	API					API						    `koanf:"api"`
	SETTINGS      		SETTINGS					`koanf:"settings"        token>aliases:"overrides" token>onuse:".overrides>>deprecated"       deprecation:"{b,s,fg=orange}\x60overrides\x60{/} is no longer needed in {b}Token Configs{/}\nUse {b,fg=green}\x60settings\x60{/} instead"`
}																														

type SERVICE struct {
	HOSTNAMES			[]string					`koanf:"hostnames"       env>aliases:".hostnames"`
	PORT				string						`koanf:"port"            env>aliases:".port"`
	LOG_LEVEL			string						`koanf:"loglevel"        env>aliases:".loglevel"`
}

type API struct {
	URL					string						`koanf:"url"             env>aliases:".apiurl"`
	TOKENS				[]string					`koanf:"tokens"          env>aliases:".apitokens,.apitoken" aliases:"token" token>aliases:".tokens,.token" token>onuse:".tokens,.token,token>>deprecated" onuse:"token>>deprecated" deprecation:".tokens,.token>>{b,s,fg=orange}\x60tokens\x60{/} and {b,s,fg=orange}\x60token\x60{/} will not be at {b}root{/} anymore\nUse {b,fg=green}\x60api.tokens\x60{/} instead|token>>{b,s,fg=orange}\x60api.token\x60{/} will be {u}removed{/} in favor of {b,fg=green}\x60api.tokens\x60{/}"`																					
	AUTH				AUTH						`koanf:"auth"`
}

type AUTH struct {
	METHODS				[]string					`koanf:"methods"         env>aliases:".authmethods"`
	TOKENS				[]Token						`koanf:"tokens"          aliases:"token" onuse:"token>>deprecated" deprecation:"{b,s,fg=orange}\x60api.auth.token\x60{/} will be removed\nUse {b,fg=green}\x60api.auth.tokens\x60{/} instead"`
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
	VARIABLES         	map[string]any              `koanf:"variables"       childtransform:"upper"`
	FIELD_MAPPINGS      map[string][]FieldMapping	`koanf:"fieldmappings"   childtransform:"default"`
	TEMPLATE  			string                      `koanf:"template"`
}

type FieldMapping struct {
	Field 				string 						`koanf:"field"`
	Score 				int    						`koanf:"score"`
}

type ACCESS struct {
	ENDPOINTS			[]string					`koanf:"endpoints"`
	FIELD_POLICIES		map[string][]FieldPolicy	`koanf:"fieldpolicies"   childtransform:"default"`
	RATE_LIMITING		RateLimiting				`koanf:"ratelimiting"`
	IP_FILTER			[]string					`koanf:"ipfilter"`
	TRUSTED_IPS			[]string					`koanf:"trustedips"`
	TRUSTED_PROXIES		[]string					`koanf:"trustedproxies"`
}

type FieldPolicy struct {
	Value				any						    `koanf:"value"`
	Action				string						`koanf:"action"`
}

type RateLimiting struct {
	Limit				int							`koanf:"limit"`
	Period				string						`koanf:"period"`
}