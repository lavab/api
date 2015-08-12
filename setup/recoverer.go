package setup

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/getsentry/raven-go"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"

	"github.com/lavab/api/env"
)

func recoverer(c *web.C, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := middleware.GetReqID(*c)

		// Raven recoverer
		defer func() {
			var packet *raven.Packet
			switch rval := recover().(type) {
			case nil:
				return
			case error:
				packet = raven.NewPacket(rval.Error(), raven.NewHttp(r), raven.NewException(rval, raven.NewStacktrace(2, 3, nil)))
			default:
				str := fmt.Sprintf("%+v", rval)
				packet = raven.NewPacket(str, raven.NewHttp(r), raven.NewException(errors.New(str), raven.NewStacktrace(2, 3, nil)))
			}

			debug.PrintStack()
			http.Error(w, http.StatusText(500), 500)

			env.Raven.Capture(packet, map[string]string{
				"request_id": id,
			})
		}()

		h.ServeHTTP(w, r)
	})
}
