package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	l "github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/gotl/pkg/request"
)

func main() {
    logLevel := os.Getenv("LOG_LEVEL")

    if strings.TrimSpace(logLevel) == "" {
        logLevel = "info"
    }

    port := os.Getenv("PORT")

    if strings.TrimSpace(port) == "" {
        port = "8881"
    }

    logger, err := l.NewWithDefaults(logLevel)

    http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
        body, err := request.GetReqBody(req)

        if err != nil {
            http.Error(w, "Bad Request: invalid body", http.StatusBadRequest)
            return
        }

        ip, _, _ := net.SplitHostPort(req.RemoteAddr)

        if body.Empty {
            logger.Info(ip, " ", req.Method, " ", req.URL.Path, " ", req.URL.RawQuery)
        } else {
            logger.Info(ip, " ", req.Method, " ", req.URL.Path, " ", req.URL.RawQuery, body.Raw)
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)

        fmt.Fprint(w, `{"message":"Hello from mock endpoint"}`)
    })

    logger.Info("Mock server running at http://127.0.0.1:", port)

    err = http.ListenAndServe("127.0.0.1:" + port, nil)

    if err != nil {
        logger.Fatal("Error starting Mock server: ", err.Error())
    }
}
