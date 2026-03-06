package bidhandler

import (
	"apps_minimalist/bidder/code/auction"

	"emperror.dev/errors"
	"github.com/valyala/fasthttp"
)

func (h Handler) writeResponse(
	response []byte,
	ctx *fasthttp.RequestCtx,
) {
	ctx.SetContentType("application/json")
	ctx.SetBody(response)
}

func (h Handler) writeResponseError(
	auctionErr error,
	ctx *fasthttp.RequestCtx,
) {
	switch {
	case errors.Is(auctionErr, auction.ErrNoBid):
		ctx.SetStatusCode(fasthttp.StatusNoContent)
	case errors.Is(auctionErr, auction.ErrTimeout):
		ctx.SetStatusCode(h.cfg.TimeoutStatus)
	default:
		ctx.Error(auctionErr.Error(), fasthttp.StatusInternalServerError)
	}
}