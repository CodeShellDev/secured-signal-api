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
	SETTINGS      		map[string]*SETTINGS
	INSECURE      		bool
}

type MESSAGE_SETTINGS struct {
	VARIABLES         	map[string]any              `koanf:"variables"`
	FIELD_MAPPINGS      map[string][]FieldMapping	`koanf:"fieldMappings"`
	TEMPLATE  			string                      `koanf:"template"`
}

type FieldMapping struct {
	Field 				string 						`koanf:"field"`
	Score 				int    						`koanf:"score"`
}

type ACCESS_SETTINGS struct {
	ENDPOINTS			[]string					`koanf:"endpoints"`
}

type SETTINGS struct {
	ACCESS 				ACCESS_SETTINGS 			`koanf:"access"`
	MESSAGE				MESSAGE_SETTINGS			`koanf:"message"`
}