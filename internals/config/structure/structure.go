package structure

import (
	t "github.com/codeshelldev/gotl/pkg/configutils/types"
	c "github.com/codeshelldev/secured-signal-api/internals/config/structure/custom"
	g "github.com/codeshelldev/secured-signal-api/internals/config/structure/generics"
)

type ENV struct {
	CONFIG_PATH   		string
	TOKENS_DIR    		string

	DEFAULTS_PATH 		string
	FAVICON_PATH  		string

	DB_PATH				string
	
	INSECURE      		bool
	REDACT_TOKENS		bool

	TOKENS				[]string

	CONFIGS      		map[string]*CONFIG
}

type CONFIG struct {
	TYPE				ConfigType
	NAME				string						`koanf:"name"`
	SERVICE				SERVICE 					`koanf:"service"`
	API					API						    `koanf:"api"`
	SETTINGS      		SETTINGS					`koanf:"settings"`
}

type ConfigType string

const (
	TOKEN ConfigType = "token"
	MAIN ConfigType = "main"
)

type SERVICE struct {
	HOSTNAMES			t.Opt[[]string]				`koanf:"hostnames"          env>aliases:".hostnames"`
	PORT				string						`koanf:"port"               env>aliases:".port"`
	LOG_LEVEL			string						`koanf:"loglevel"           env>aliases:".loglevel"`
}

type API struct {
	URL					*g.URL						`koanf:"url"                env>aliases:".apiurl"`
	TOKENS				[]string					`koanf:"tokens"             env>aliases:".apitokens"`
	AUTH				AUTH						`koanf:"auth"`
}

type AUTH struct {
	METHODS				t.Opt[[]string]				`koanf:"methods"            env>aliases:".authmethods"`
	TOKENS				[]Token						`koanf:"tokens"`
}

type Token struct {
	Set					[]string					`koanf:"set"`
	Methods				[]string					`koanf:"methods"`
}

type SETTINGS struct {
	ACCESS 				ACCESS 						`koanf:"access"`
	MESSAGE				MESSAGE						`koanf:"message"`
	HTTP				HTTP						`koanf:"http"`
}

type HTTP struct {
	RESPONSE_HEADERS	t.Opt[map[string]string]	`koanf:"responseheaders"`
}

type MESSAGE struct {
	VARIABLES         	t.Opt[map[string]any]		`koanf:"variables"          childtransform:"upper"`
	FIELD_MAPPINGS      t.Opt[FieldMappings]		`koanf:"fieldmappings"      childtransform:"default"`
	TEMPLATING  		t.Opt[Templating]			`koanf:"templating"         aliases:"template"          onuse:"template>>broken"    breaking:"{b,fg=red}\x60{s}settings.message.template{/}\x60{/} has been moved\n Use {b,fg=green}\x60settings.message.templating.messageTemplate\x60{/} instead"`
	SCHEDULING			t.Opt[Scheduling]			`koanf:"scheduling"`
	INJECTING			t.Opt[Injecting]			`koanf:"injecting"`
}

type FieldMappings = map[string][]FMapping

type Injecting struct {
	URLToBody			t.Opt[URLToBody]			`koanf:"urltobody"`
}

type URLToBody struct {
	Path				bool						`koanf:"path"`
	Query				bool						`koanf:"query"`
}

type Templating struct {
	MessageTemplate		string						`koanf:"messagetemplate"`
	Body				bool						`koanf:"body"`
	Query				bool						`koanf:"query"`
	Path				bool						`koanf:"path"`
}

type Scheduling struct {
	// Enabled is needed because this isn't a data-driven setting, but rather a toggle
	Enabled				bool						`koanf:"enabled"`
	MaxHorizon			t.Opt[g.TimeDuration]		`koanf:"maxhorizon"`
}

type FMapping struct {
	Field 				string 						`koanf:"field"`
	Score 				int    						`koanf:"score"`
}

type ACCESS struct {
	ENDPOINTS			t.Opt[Endpoints] 			`koanf:"endpoints"`
	FIELD_POLICIES		t.Opt[FieldPolicies] 		`koanf:"fieldpolicies"      childtransform:"default"`
	RATE_LIMITING		t.Opt[RateLimiting]			`koanf:"ratelimiting"`
	IP_FILTER			t.Opt[IPFilter]				`koanf:"ipfilter"`
	TRUSTED_IPS			t.Opt[[]g.IPOrNet]			`koanf:"trustedips"`
	TRUSTED_PROXIES		t.Opt[[]g.IPOrNet]			`koanf:"trustedproxies"`
	CORS				t.Opt[Cors]					`koanf:"cors"`
}

type Cors struct {
	Origins				[]Origin					`koanf:"origins"`
	Methods				t.Opt[[]string]				`koanf:"methods"`
	Headers				t.Opt[[]string]				`koanf:"headers"`
}

type Origin struct {
	URL					g.URL						`koanf:"url"`
	Methods				t.Opt[[]string]				`koanf:"methods"`
	Headers				t.Opt[[]string]				`koanf:"headers"`
}

type FieldPolicies = *t.Comp[c.RFieldPolicies, c.FieldPolicies]

type Endpoints struct {
	Allowed				[]g.StringMatchRule			`koanf:"allowed"`
	Blocked				[]g.StringMatchRule			`koanf:"blocked"`
}

type IPFilter struct {
	Allowed				[]g.IPOrNet					`koanf:"allowed"`
	Blocked				[]g.IPOrNet					`koanf:"blocked"`
}

type FPolicy struct {
	Match				g.MatchRule[any]			`koanf:"match"`
	Action				string						`koanf:"action"`
}

type RateLimiting struct {
	Limit				int							`koanf:"limit"`
	Period				g.TimeDuration				`koanf:"period"`
}