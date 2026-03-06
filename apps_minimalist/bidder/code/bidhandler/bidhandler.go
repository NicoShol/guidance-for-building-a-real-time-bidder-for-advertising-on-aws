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

		// byteRequest, request := h.readRequest(ctx, pd)  // byteRequest will be snet to datastream ? 
		_, request := h.readRequest(ctx, pd)
		if request == nil {
			return
		}

		response := auction.Response{}
		err := h.auction.Run(deadline, request, &response)

		log.Info().Msg("Received request")
		if err != nil {
			h.writeResponseError(err, ctx)
			return
		}

		bidResponse := buildResponse(&response, pd)
		h.writeResponse(bidResponse, ctx)
	}