package middlewares

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/gotl/pkg/request"
	"github.com/codeshelldev/secured-signal-api/internals/config"
)

var InternalAPI Middleware = Middleware{
	Name: "_Internal_API",
	Use: internalAPIHandler,
}

func internalAPIHandler(next http.Handler) http.Handler {
	mux := http.NewServeMux()

	const aboutEndpoint = "/v1/about"
	mux.HandleFunc(aboutEndpoint, func(w http.ResponseWriter, req *http.Request) {
		res, err := redirectRequestToBackend(req, config.DEFAULT.API.URL + aboutEndpoint)
		
		if err != nil {
			logger.Error("Error requesting Backend: ", err.Error())
			http.Error(w, "Internal Server Error ", http.StatusBadGateway)
			return
		}

		body, err := request.GetResBody(res)

		if err != nil {
			logger.Error("Could not get Response Body: ", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		for key, values := range res.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		if !body.Empty {
			var version string

			if isValidSemver(os.Getenv("IMAGE_TAG")) {
				version, _ = strings.CutPrefix(version, "v")
			}
			
			payload := map[string]any{
				"version": version,
				"auth_required": !config.ENV.INSECURE,
			}

			body.Data["secured-signal-api"] = payload

			err := body.Write(w)

			if err != nil {
				logger.Error("Could not write to Response Body: ", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(res.StatusCode)
	})

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, config.ENV.FAVICON_PATH)
	})

	mux.Handle("/", next)

	return mux
}

func isValidSemver(version string) bool {
	re, err := regexp.Compile(`^v?([0-9]+)\.([0-9]+)\.([0-9]+)(-[0-9A-Za-z.-]+)?(\+[0-9A-Za-z.-]+)?$`)
	
	if err != nil {
		return false
	}

	return re.MatchString(version)
}

func redirectRequestToBackend(req *http.Request, url string) (*http.Response, error) {
	var body io.Reader

    if req.Body != nil {
        bodyBytes, err := io.ReadAll(req.Body)

        if err != nil {
            return nil, err
        }
        req.Body.Close()

        body = bytes.NewReader(bodyBytes)
        req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
    }

	return requestBackend(req.Method, url, body, req.Header)
}

func requestBackend(method, url string, body io.Reader, headers map[string][]string) (*http.Response, error) {
    backendReq, err := http.NewRequest(method, url, body)
    if err != nil {
        return nil, err
    }

    for key, values := range headers {
        for _, value := range values {
            backendReq.Header.Add(key, value)
        }
    }

    client := &http.Client{}

    res, err := client.Do(backendReq)
    if err != nil {
        return nil, err
    }

    return res, nil
}