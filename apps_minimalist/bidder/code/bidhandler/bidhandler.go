package bidhandler

import (
	"apps_minimalist/bidder/code/auction"
	"apps_minimalist/bidder/code/stream"

	"time"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

type Handler struct {
	cfg       	*Config
	auction   	*auction.Auction
	pool 		*pool
	stream     	stream.Stream
}

func New(
	cfg Config,
	auctionFn *auction.Auction,
	dataStream stream.Stream,
) Handler {
	return Handler{
		cfg: &cfg,
		auction: auctionFn,
		pool: &pool{},
		stream: dataStream,
	}
}


func (h Handler) HandleRequest(ctx *fasthttp.RequestCtx) {
		log.Info().Msg("Recieved request")
		deadline := time.Now().Add(h.cfg.Timeout)

		// Get persistent data from the pool
		pd := h.pool.Get()
		defer h.pool.Put(pd)

		byteRequest, request := h.readRequest(ctx, pd)
		if request == nil { return }
		if h.stream != nil {
			h.stream.PutRequest(byteRequest)
		}

		response := auction.Response{}
		err := h.auction.Run(deadline, request, &response)

		log.Info().Msg("Received request")
		if err != nil {
			h.writeResponseError(err, ctx)
			return
		}

		bidResponse := buildResponse(&response, pd)
		if h.stream != nil {
			h.stream.PutResponse(bidResponse)
		}
		h.writeResponse(bidResponse, ctx)
	}