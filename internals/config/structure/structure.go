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
	SETTINGS      		SETTINGS					`koanf:"settings"        token>aliases:"overrides" token>onuse:".overrides>>deprecated"       deprecation:"We have moved away from using 'overrides' in Token Configs"`
}																														

type SERVICE struct {
	HOSTNAMES			[]string					`koanf:"hostnames"       env>aliases:".hostnames"`
	PORT				string						`koanf:"port"            env>aliases:".port"`
	LOG_LEVEL			string						`koanf:"loglevel"        env>aliases:".loglevel"`
}

type API struct {
	URL					string						`koanf:"url"             env>aliases:".apiurl"`
	TOKENS				[]string					`koanf:"tokens"          env>aliases:".apitokens,.apitoken" aliases:"token" token>aliases:".tokens,.token" token>onuse:".tokens>>deprecated,.token>>deprecated,token>>deprecated" onuse:"token>>deprecated" deprecation:"'tokens' and 'token' will not be at the root anymore\n'api.token' will be removed in favor of 'api.tokens'"`																					
	AUTH				AUTH						`koanf:"auth"`
}

type AUTH struct {
	METHODS				[]string					`koanf:"methods"         env>aliases:".authmethods"`
	TOKENS				[]Token						`koanf:"tokens"          aliases:"token" onuse:"token>>deprecated" deprecation:"'api.auth.token' will be removed in favor of 'api.auth.tokens'"`
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