package structure

type ENV struct {
	CONFIG_PATH   		string
	DEFAULTS_PATH 		string
	FAVICON_PATH  		string
	TOKENS_DIR    		string
	LOG_LEVEL     		string
	PORT          		string
	API_URL       		string
	API_TOKENS    		[]string
	SETTINGS      		map[string]*SETTINGS		`koanf:"settings"`
	INSECURE      		bool
}

type SETTINGS struct {
	ACCESS 				ACCESS_SETTINGS 			`koanf:"access"        transform:"lower"`
	MESSAGE				MESSAGE_SETTINGS			`koanf:"message"       transform:"lower"`
}

type MESSAGE_SETTINGS struct {
	VARIABLES         	map[string]any              `koanf:"variables"                       childtransform:"upper"`
	FIELD_MAPPINGS      map[string][]FieldMapping	`koanf:"fieldmappings"                   childtransform:"default"`
	TEMPLATE  			string                      `koanf:"template"      transform:"lower"`
}

type FieldMapping struct {
	Field 				string 						`koanf:"field"         transform:"lower"`
	Score 				int    						`koanf:"score"         transform:"lower"`
}

type ACCESS_SETTINGS struct {
	ENDPOINTS			[]string					`koanf:"endpoints"     transform:"lower"`
	FIELD_POLICIES		map[string]FieldPolicy		`koanf:"fieldpolicies" transform:"lower" childtransform:"default"`
}

type FieldPolicy struct {
	Value				any						    `koanf:"value"         transform:"lower"`
	Action				string						`koanf:"action"        transform:"lower"`
}