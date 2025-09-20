package middlewares

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	middlewareTypes "github.com/codeshelldev/secured-signal-api/internals/proxy/middlewares/types"
	jsonutils "github.com/codeshelldev/secured-signal-api/utils/jsonutils"
	log "github.com/codeshelldev/secured-signal-api/utils/logger"
	request "github.com/codeshelldev/secured-signal-api/utils/request"
)

type AliasMiddleware struct {
	Next 	http.Handler
}

func (data AliasMiddleware) Use() http.Handler {
	next := data.Next

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		settings := getSettingsByReq(req)

		dataAliases := settings.DATA_ALIASES

		if dataAliases == nil {
			dataAliases = getSettings("*").DATA_ALIASES
		}

		if settings.VARIABLES == nil {
			settings.VARIABLES = getSettings("*").VARIABLES
		}

		body, err := request.GetReqBody(w, req)

		if err != nil {
			log.Error("Could not get Request Body: ", err.Error())
		}

		var modifiedBody bool
		var bodyData map[string]any

		if !body.Empty {
			bodyData = body.Data

			aliasData := processDataAliases(dataAliases, bodyData)

			for key, value := range aliasData {
				prefix := key[:1]

				keyWithoutPrefix := key[1:]

				switch prefix {
					case "@":
						bodyData[keyWithoutPrefix] = value
						modifiedBody = true
					case ".":
						settings.VARIABLES[keyWithoutPrefix] = value
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

func processDataAliases(aliases map[string][]middlewareTypes.DataAlias, data map[string]any) (map[string]any) {
	aliasData := map[string]any{}

	for key, alias := range aliases {
		key, value := getData(key, alias, data)

		aliasData[key] = value
	}

	return aliasData
}

func getData(key string, aliases []middlewareTypes.DataAlias, data map[string]any) (string, any) {
	var best int
	var value any

	for _, alias := range aliases {
		aliasValue, score, ok := processAlias(alias, data)

		if ok {
			if score > best {
				value = aliasValue
			}

			delete(data, alias.Alias)
		}
	}

	return key, value
}

func processAlias(alias middlewareTypes.DataAlias, data map[string]any) (any, int, bool) {
	aliasKey := alias.Alias

	value, ok := jsonutils.GetByPath(aliasKey, data)

	if ok && value != nil {
		return value, alias.Score, true
	} else {
		return "", 0, false
	}
}