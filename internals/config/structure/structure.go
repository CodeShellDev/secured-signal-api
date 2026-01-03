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
																			//TODO: deprecate overrides for tkconfigs
	SETTINGS      		SETTINGS					`koanf:"settings"        token>aliases:"overrides"`
}

type SERVICE struct {
	PORT				string						`koanf:"port"            env>aliases:".port"`
	LOG_LEVEL			string						`koanf:"loglevel"        env>aliases:".loglevel"`
}

type API struct {
	URL					string						`koanf:"url"             env>aliases:".apiurl"`
																													//TODO: deprecate .token for tkconfigs
	TOKENS				[]string					`koanf:"tokens"          env>aliases:".apitokens,.apitoken"     token>aliases:".tokens,.token"       aliases:"token"`
}

type SETTINGS struct {
	ACCESS 				ACCESS 			`koanf:"access"`
	MESSAGE				MESSAGE			`koanf:"message"`
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
	FIELD_POLICIES		map[string]FieldPolicy		`koanf:"fieldpolicies"   childtransform:"default"`
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