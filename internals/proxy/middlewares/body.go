package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	log "github.com/codeshelldev/secured-signal-api/utils/logger"
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
		if req.Body != nil {
			bodyBytes, err := io.ReadAll(req.Body)

			if err != nil {
				log.Error("Could not read Body: ", err.Error())
				http.Error(w, "Internal Error", http.StatusInternalServerError)
				return
			}

			req.Body.Close()

			var modifiedBodyData map[string]interface{}

			err = json.Unmarshal(bodyBytes, &modifiedBodyData)

			if err != nil {
				log.Error("Could not decode Body: ", err.Error())
				http.Error(w, "Internal Error", http.StatusInternalServerError)
				return
			}

			content, ok := modifiedBodyData["message"]

			if !ok || content == "" {
				best := 0

				for _, alias := range messageAliases {
					aliasKey := alias.Alias
					priority := alias.Priority

					value, ok := modifiedBodyData[aliasKey]

					if ok && value != "" && priority > best {
						content = modifiedBodyData[aliasKey]
					}

					modifiedBodyData[aliasKey] = nil
				}

				modifiedBodyData["message"] = content

				bodyBytes, err = json.Marshal(modifiedBodyData)

				if err != nil {
					log.Error("Could not encode Body: ", err.Error())
					http.Error(w, "Internal Error", http.StatusInternalServerError)
					return
				}
			}

			modifiedBody := string(bodyBytes)

			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))

			req.ContentLength = int64(len(modifiedBody))
			req.Header.Set("Content-Length", strconv.Itoa(len(modifiedBody)))
		}

		next.ServeHTTP(w, req)
	})
}
