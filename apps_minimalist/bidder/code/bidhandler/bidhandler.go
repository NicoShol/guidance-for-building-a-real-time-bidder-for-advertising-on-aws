package bidhandler

import (
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

type Handler struct {
	cfg        *Config
}

func New(
	cfg Config,
) Handler {
	return Handler{
		cfg: &cfg,
	}
}

func (h Handler) HandleRequest(ctx *fasthttp.RequestCtx) {
		// // TESTING HERE : simply log requests
		log.Info().
			Str("method", string(ctx.Method())).
			Str("path", string(ctx.Path())).
			Msg("recieved request")

		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetContentType("text/plain")
		ctx.SetBodyString("Hello from bidder!")
		// // END TESTING
}