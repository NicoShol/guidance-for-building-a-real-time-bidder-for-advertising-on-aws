package bidhandler

import (
	"apps_minimalist/bidder/code/auction"
	"apps_minimalist/bidder/code/openrtb"

	"emperror.dev/errors"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

var errBadVersion = errors.New("unsupported x-openrtb-version")

func (h Handler) readRequest(
	ctx *fasthttp.RequestCtx,
	pd *persistentData,
) ([]byte, *auction.ExtendedRequest) {
		byteRequest := ctx.PostBody()
		openRTBVersion := openrtb.Version(ctx.Request.Header.Peek("x-openrtb-version"))
		
		if openRTBVersion == "" {
			openRTBVersion = h.cfg.OpenRTBVersion
		}

		var err error
		var request *auction.ExtendedRequest

		switch openRTBVersion {
		case openrtb.OpenRTB2_5:
			request, err = parseBidRequestExtended_2_5(byteRequest, pd)
		default:
			err = errBadVersion
		}

		if err != nil {
			log.Error().Err(err).Msg("")
			ctx.Error(err.Error(), fasthttp.StatusBadRequest)
			return byteRequest, nil
		}
		return byteRequest, request
}

// func parseBidRequest2(byteRequest []byte, pd *persistentData) (*auction.Request, error) {
// 	v, err := pd.parser.ParseBytes(byteRequest)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "error while parsing request")
// 	}

// 	ID := v.GetStringBytes("id")

// 	// TODO : continue
// }
