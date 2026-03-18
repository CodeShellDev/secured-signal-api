package middleware

import (
	"net/http"

	. "github.com/codeshelldev/secured-signal-api/internals/proxy/common"
)

var InternalResponseHooks = ResponseMiddleware{
	Name: "_Response_Hooks",
	Use: hooksHandler,
}

func hooksHandler(res *http.Response) error {
	hooks := GetResponseHooks(res.Request)

	for _, h := range hooks {
		err := h(res)

		if err != nil {
			return err
		}
	}

	return nil
}