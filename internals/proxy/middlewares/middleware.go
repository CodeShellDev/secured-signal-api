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