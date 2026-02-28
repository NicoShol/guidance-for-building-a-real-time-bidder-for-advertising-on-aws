package server

import (
	// "github.com/rs/zerolog/log"
	"apps_minimalist/bidder/code/bidhandler"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
)

type Server struct {
	server fasthttp.Server
	cfg Config
}

func NewServer(
	cfg Config,
	bidHandler bidhandler.Handler,
) *Server {
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case cfg.BidRequestPath:
			if !ctx.IsPost() {
				ctx.Error("Unsupported method", fasthttp.StatusMethodNotAllowed)
				return
			} 
			bidHandler.HandleRequest(ctx)
		default:
			ctx.Error("Unsupported path", fasthttp.StatusNotFound)
		}
	}

	return &Server{
		server: fasthttp.Server{
			Handler: 	requestHandler,
			Logger:		logger{},
			// Other server configurations can be set here
			NoDefaultServerHeader: true,
		},
		cfg: cfg,
	}
}

func (s *Server) AsyncListenAndServe(errCallback func(error)) {
	// go func() { ... } () 
	// go : keyword for "run this concurently in the background"
	// func() {...} anonymous func definition 
	// () at the end => invoke immediatly
	go func() {
		listener, err := reuseport.Listen("tcp4", s.cfg.Address)
		if err != nil && errCallback != nil {
			errCallback(err)
			return
		}
		if err := s.server.Serve(listener); err != nil {
			if errCallback != nil {
				errCallback(err)
			}
		}
	} ()
}