package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/codeshelldev/gotl/pkg/logger"
)

type Server struct {
	Host string
	Ports []string
	Handler http.Handler
	Listeners map[string]*http.Server
}

func Create(handler http.Handler, host string, ports ...string) *Server {
	return &Server{
		Host: host,
		Ports: ports,
		Handler: handler,
		Listeners: map[string]*http.Server{},
	}
}

func (server *Server) ListenAndServer() {
    var wg sync.WaitGroup
    stopCh := make(chan struct{})

	for _, port := range server.Ports {
		addr := server.Host + ":" + port
		listener, err := net.Listen("tcp", addr)

		if err != nil {
			logger.Error("Error listening on ", port, ": ", err.Error())
			continue
		}

		srv := &http.Server{
			Addr: server.Host + ":" + port,
			Handler: server.Handler,
		}

		wg.Add(1)

		go func(s *http.Server, l net.Listener, p string) {
            defer wg.Done()

			logger.Debug("Listener on port ", port, " started")

			server.Listeners[port] = s

            err := s.Serve(l)

            if err != nil && err != http.ErrServerClosed {
				logger.Error("Listener on port ", port, " exited with ", err.Error())
            }
        }(srv, listener, port)
	}

	go func() {
        wg.Wait()
        close(stopCh)
    }()

	<- stopCh
}

func (server *Server) Shutdown(ctx context.Context) error {
	var errs []error

	for port, s := range server.Listeners {
		logger.Debug("Shutting down listener on ", port)

		err := s.Shutdown(ctx)

		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func PortsToRangeString(ports []string) string {
    if len(ports) == 0 {
        return ""
    }

    sort.Strings(ports)

    result := []string{}

	end, _ := strconv.Atoi(ports[0])
	start, _ := strconv.Atoi(ports[0])

    for i := 1; i < len(ports); i++ {
		port, _ := strconv.Atoi(ports[i])

        if port == end + 1 {
            end = port
        } else {
            if start == end {
                result = append(result, strconv.Itoa(start))
            } else {
                result = append(result, strconv.Itoa(start) + "-" + strconv.Itoa(end))
            }

            start = port
            end = port
        }
    }

    if start == end {
        result = append(result, strconv.Itoa(start))
    } else {
		result = append(result, strconv.Itoa(start) + "-" + strconv.Itoa(end))
    }

    return strings.Join(result, ",")
}