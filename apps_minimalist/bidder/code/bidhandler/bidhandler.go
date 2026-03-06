package bidhandler

import (
	"apps_minimalist/bidder/code/auction"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

type Handler struct {
	cfg        *Config
	auction    *auction.Auction
	pool 		*pool
}

func New(
	cfg Config,
	auctionFn *auction.Auction,
) Handler {
	return Handler{
		cfg: &cfg,
		auction: auctionFn,
		pool: &pool{},
	}
}

func (h Handler) HandleRequest(ctx *fasthttp.RequestCtx) {
		log.Info().Msg("Recieved request")
		deadline := time.Now().Add(h.cfg.Timeout)

		// Get persistent data from the pool
		pd := h.pool.Get()
		defer h.pool.Put(pd)

		byteRequest, request := h.readRequest(ctx, pd)
		if request == nil {
			return
		}

		response := auction.Response{}
		err := h.auction.Run(deadline, request, &response)
		// TMP error handling, to be replaced with proper response writing in case of error
		if err != nil {
			log.Error().Err(err).Msg("Error while running auction")
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}
		// TODO: handle error and set response in context 
		// for instance: 
		// log.Info().Msg("Received request")
		// if err != nil {
		// 	h.writeResponseError(err, ctx)
		// 	return
		// }

		// bidResponse := buildResponse(&response, pd)
		// h.writeResponse(bidResponse, ctx)
		// TESTING -- just log the byte request for now
		log.Info().Bytes("request", byteRequest).Msg("Parsed request")
		// END TESTING

		// log.Info().Msg("Parsed request")
		// response := auction.Response{}


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