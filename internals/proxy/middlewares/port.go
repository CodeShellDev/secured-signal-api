package middlewares

import (
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/codeshelldev/gotl/pkg/logger"
)

var Port Middleware = Middleware{
	Name: "Port",
	Use: portHandler,
}

func portHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		conf := getConfigByReq(req)

		allowedPort := conf.SERVICE.PORT

		if strings.TrimSpace(allowedPort) == "" {
			next.ServeHTTP(w, req)
			return
		}

		port, err := getPort(req)

		if err != nil {
			logger.Error("Could not get Port: ", err.Error())
			http.Error(w, "Bad Request: invalid scheme", http.StatusBadRequest)
			return
		}

		if port != allowedPort {
			logger.Warn("User tried using Token on wrong Port")
			onUnauthorized(w)
			return
		}
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