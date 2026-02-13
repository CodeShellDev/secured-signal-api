package endpoints

import (
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/gotl/pkg/request"
	"github.com/codeshelldev/secured-signal-api/internals/config"
	. "github.com/codeshelldev/secured-signal-api/internals/proxy/common"
)

var AboutEndpoint = Endpoint{
	Name: "About",
	Handler: aboutHandler,
}

func aboutHandler(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("GET /v1/about", func(w http.ResponseWriter, req *http.Request) {
		req.RequestURI = ""
		ChangeRequestDest(req, config.DEFAULT.API.URL.String() + "/v1/about")

		conf := GetConfigByReq(req)

		client := &http.Client{}
		res, err := client.Do(req)

		if err != nil {
			logger.Error("Error requesting backend: ", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
				"capabilities": map[string]any{
					"v2/send": getSendCapabilities(conf),
				},
			}

			body.Data["secured-signal-api"] = payload

			err := body.Write(w)

			if err != nil {
				logger.Error("Could not write to Response Body: ", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
	})

	return mux
}

func isValidSemver(version string) bool {
	re, err := regexp.Compile(`^v?([0-9]+)\.([0-9]+)\.([0-9]+)(-[0-9A-Za-z.-]+)?(\+[0-9A-Za-z.-]+)?$`)
	
	if err != nil {
		return false
	}

	return re.MatchString(version)
}