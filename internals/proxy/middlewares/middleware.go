package middlewares

import (
	"net/http"
	"strings"

	"github.com/codeshelldev/gotl/pkg/logger"
)

type Middleware struct {
	Name string
	Use func(next http.Handler) http.Handler 
}

type Chain struct {
    middlewares []Middleware
}

func NewChain() *Chain {
    return &Chain{}
}

func (chain *Chain) Use(middleware Middleware) *Chain {
    chain.middlewares = append(chain.middlewares, middleware)

    if strings.HasPrefix(middleware.Name, "_") {
        logger.Dev("Registered ", middleware.Name, " middleware")
    } else {
	    logger.Debug("Registered ", middleware.Name, " middleware")
    }


    return chain
}

func (chain *Chain) Then(final http.Handler) http.Handler {
    handler := final

    for i := len(chain.middlewares) - 1; i >= 0; i-- {
        handler = chain.middlewares[i].Use(handler)
    }

    return handler
}

type ResponseMiddleware struct {
	Name string
	Use func(res *http.Response) error 
}

type ResponseChain struct {
    middlewares []ResponseMiddleware
}

func NewResponseChain() *ResponseChain {
    return &ResponseChain{}
}

func (chain *ResponseChain) Use(middleware ResponseMiddleware) *ResponseChain {
    chain.middlewares = append(chain.middlewares, middleware)

    if strings.HasPrefix(middleware.Name, "_") {
        logger.Dev("Registered ", middleware.Name, " response middleware")
    } else {
	    logger.Debug("Registered ", middleware.Name, " response middleware")
    }

    return chain
}

func (chain *ResponseChain) Then() func(*http.Response) error {
    return func(resp *http.Response) error {
        for _, middleware := range chain.middlewares {
            err := middleware.Use(resp)

            if err != nil {
                return err
            }
        }

        return nil
    }
}