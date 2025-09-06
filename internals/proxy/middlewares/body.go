package middlewares

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	request "github.com/codeshelldev/secured-signal-api/utils/request"
)

type MessageAlias struct {
	Alias    string
	Priority int
}

type BodyMiddleware struct {
	Next           http.Handler
	MessageAliases []MessageAlias
}

func (data BodyMiddleware) Use() http.Handler {
	next := data.Next
	messageAliases := data.MessageAliases

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		body := request.GetReqBody(w, req)

		var modifiedBody bool
		var bodyData map[string]interface{}

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

func getMessage(aliases []MessageAlias, data map[string]interface{}) (string, map[string]interface{}) {
	var content string
	var best int

	for _, alias := range aliases {
		aliasKey := alias.Alias
		priority := alias.Priority

		value, ok := data[aliasKey]

		if ok && value != "" && priority > best {
			content = data[aliasKey].(string)
		}

		data[aliasKey] = nil
	}

	return content, data
}