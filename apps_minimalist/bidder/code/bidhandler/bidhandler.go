package bidhandler

import (
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

type Handler struct {
	cfg        *Config

	pool 		*pool
}

func New(
	cfg Config,
) Handler {
	return Handler{
		cfg: &cfg,
		pool: &pool{},
	}
}

func (h Handler) HandleRequest(ctx *fasthttp.RequestCtx) {
		log.Info().Msg("Recieved request")

		// Get persistent data from the pool
		pd := h.pool.Get()
		defer h.pool.Put(pd)

		byteRequest, request := h.readRequest(ctx, pd)
		if request == nil {
			return
		}
		// // TESTING HERE : simply log requests
		// log.Info().
		// 	Str("method", string(ctx.Method())).
		// 	Str("path", string(ctx.Path())).
		// 	Msg("recieved request")

		// ctx.SetStatusCode(fasthttp.StatusOK)
		// ctx.SetContentType("text/plain")
		// ctx.SetBodyString("Hello from bidder!")
		// // END TESTING
}