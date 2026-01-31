package middlewares

import (
	"net/http"

	jsonutils "github.com/codeshelldev/gotl/pkg/jsonutils"
	request "github.com/codeshelldev/gotl/pkg/request"
	"github.com/codeshelldev/secured-signal-api/internals/config/structure"
)

var Mapping Middleware = Middleware{
	Name: "Mapping",
	Use: mappingHandler,
}

func mappingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger := getLogger(req)

		conf := getConfigByReq(req)

		variables := conf.SETTINGS.MESSAGE.VARIABLES
		fieldMappings := conf.SETTINGS.MESSAGE.FIELD_MAPPINGS

		if fieldMappings == nil {
			fieldMappings = getConfig("").SETTINGS.MESSAGE.FIELD_MAPPINGS
		}

		if variables == nil {
			variables = getConfig("").SETTINGS.MESSAGE.VARIABLES
		}

		body, err := request.GetReqBody(req)

		if err != nil {
			logger.Error("Could not get Request Body: ", err.Error())
			http.Error(w, "Bad Request: invalid body", http.StatusBadRequest)
			return
		}

		var modifiedBody bool
		var bodyData map[string]any

		if !body.Empty {
			bodyData = body.Data

			aliasData := processFieldMappings(fieldMappings, bodyData)

			for key, value := range aliasData {
				prefix := key[:1]

				keyWithoutPrefix := key[1:]

				switch prefix {
				case "@":
					bodyData[keyWithoutPrefix] = value
					modifiedBody = true
				case ".":
					variables[keyWithoutPrefix] = value
				}
			}
		}

		if modifiedBody {
			body.Data = bodyData

			err := body.UpdateReq(req)

			if err != nil {
				logger.Error("Could not write to Request Body: ", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			logger.Debug("Applied Data Aliasing: ", body.Data)
		}

		next.ServeHTTP(w, req)
	})
}

func processFieldMappings(aliases map[string][]structure.FieldMapping, data map[string]any) map[string]any {
	aliasData := map[string]any{}

	for key, alias := range aliases {
		key, value := getData(key, alias, data)

		if value != nil {
			aliasData[key] = value
		}
	}

	return aliasData
}

func getData(key string, aliases []structure.FieldMapping, data map[string]any) (string, any) {
	var best int
	var value any

	for _, alias := range aliases {
		aliasValue, score, ok := processFieldMapping(alias, data)

		if ok {
			if score > best {
				value = aliasValue
			}

			delete(data, alias.Field)
		}
	}

	return key, value
}

func processFieldMapping(alias structure.FieldMapping, data map[string]any) (any, int, bool) {
	aliasKey := alias.Field

	value, ok := jsonutils.GetByPath(aliasKey, data)

	if ok && value != nil {
		return value, alias.Score, true
	} else {
		return "", 0, false
	}
}
