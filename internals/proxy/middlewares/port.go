package middlewares

import (
	"errors"
	"net"
	"net/http"
	"strings"

	. "github.com/codeshelldev/secured-signal-api/internals/proxy/common"
)

var Port Middleware = Middleware{
	Name: "Port",
	Use: portHandler,
}

func portHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger := GetLogger(req)

		conf := GetConfigByReq(req)

		allowedPort := conf.SERVICE.PORT

		if strings.TrimSpace(allowedPort) == "" {
			next.ServeHTTP(w, req)
			return
		}

		port, err := getPort(req)

		if err != nil {
			logger.Error("Could not get Port: ", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if port != allowedPort {
			logger.Warn("Client tried using Token on wrong Port")
			onUnauthorized(w)
			return
		}

		next.ServeHTTP(w, req)
	})
}

func getPort(req *http.Request) (string, error) {
    addr, ok := req.Context().Value(http.LocalAddrContextKey).(net.Addr)

    if !ok {
        return "", errors.New("no local addr in context")
    }

    _, port, err := net.SplitHostPort(addr.String())

    return port, err
}