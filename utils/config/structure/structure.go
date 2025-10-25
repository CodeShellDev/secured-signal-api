package structure

type FieldMapping struct {
	Field string `koanf:"field"`
	Score int    `koanf:"score"`
}

type ENV_ struct {
	CONFIG_PATH   string
	DEFAULTS_PATH string
	FAVICON_PATH  string
	TOKENS_DIR    string
	LOG_LEVEL     string
	PORT          string
	API_URL       string
	API_TOKENS    []string
	SETTINGS      map[string]*SETTING_
	INSECURE      bool
}

type SETTING_ struct {
	ENDPOINTS 			[]string 					`koanf:"access.endpoints"`
	VARIABLES         	map[string]any              `koanf:"message.variables"`
	FIELD_MAPPINGS      map[string][]FieldMapping	`koanf:"message.fieldMappings"`
	MESSAGE_TEMPLATE  	string                      `koanf:"message.template"`
}