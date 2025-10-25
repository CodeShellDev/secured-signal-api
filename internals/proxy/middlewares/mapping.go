package middlewares

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/codeshelldev/secured-signal-api/utils/config/structure"
	jsonutils "github.com/codeshelldev/secured-signal-api/utils/jsonutils"
	log "github.com/codeshelldev/secured-signal-api/utils/logger"
	request "github.com/codeshelldev/secured-signal-api/utils/request"
)

var Mapping Middleware = Middleware{
	Name: "Mapping",
	Use: mappingHandler,
}

func mappingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		settings := getSettingsByReq(req)

		fieldMappings := settings.MESSAGE.FIELD_MAPPINGS

		if fieldMappings == nil {
			fieldMappings = getSettings("*").MESSAGE.FIELD_MAPPINGS
		}

		if settings.MESSAGE.VARIABLES == nil {
			settings.MESSAGE.VARIABLES = getSettings("*").MESSAGE.VARIABLES
		}

		body, err := request.GetReqBody(w, req)

		if err != nil {
			log.Error("Could not get Request Body: ", err.Error())
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
					settings.MESSAGE.VARIABLES[keyWithoutPrefix] = value
				}
			}
		}

		if modifiedBody {
			modifiedBody, err := request.CreateBody(bodyData)

			if err != nil {
				http.Error(w, "Internal Error", http.StatusInternalServerError)
				return
			}

			body = modifiedBody

			strData := body.ToString()

			log.Debug("Applied Data Aliasing: ", strData)

			req.ContentLength = int64(len(strData))
			req.Header.Set("Content-Length", strconv.Itoa(len(strData)))
		}

		req.Body = io.NopCloser(bytes.NewReader(body.Raw))

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
