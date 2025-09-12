package middlewares

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/codeshelldev/secured-signal-api/utils"
	log "github.com/codeshelldev/secured-signal-api/utils/logger"
	request "github.com/codeshelldev/secured-signal-api/utils/request"
)

type MessageAlias struct {
	Alias    string
	Score 	 int
}

type BodyMiddleware struct {
	Next           http.Handler
	MessageAliases []MessageAlias
}

func (data BodyMiddleware) Use() http.Handler {
	next := data.Next
	messageAliases := data.MessageAliases

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		body, err := request.GetReqBody(w, req)

		if err != nil {
			log.Error("Could not get Request Body: ", err.Error())
		}

		var modifiedBody bool
		var bodyData map[string]any

		if !body.Empty {
			bodyData = body.Data

			content, ok := bodyData["message"]

			if !ok || content == "" {

				bodyData["message"], bodyData = getMessage(messageAliases, bodyData)

				modifiedBody = true
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

			req.ContentLength = int64(len(strData))
			req.Header.Set("Content-Length", strconv.Itoa(len(strData)))
		}

		req.Body = io.NopCloser(bytes.NewReader(body.Raw))

		next.ServeHTTP(w, req)
	})
}

func getMessage(aliases []MessageAlias, data map[string]any) (string, map[string]any) {
	var content string
	var best int

	for _, alias := range aliases {
		aliasValue, score, ok := processAlias(alias, data)

		if ok && score > best {
			content = aliasValue
		}

		data[alias.Alias] = nil
	}

	return content, data
}

func processAlias(alias MessageAlias, data map[string]any) (string, int, bool) {
	aliasKey := alias.Alias

	value, ok := utils.GetByPath(aliasKey, data)

	aliasValue, isStr := value.(string)

	if isStr && ok && aliasValue != "" {
		return aliasValue, alias.Score, true
	} else {
		return "", 0, false
	}
}