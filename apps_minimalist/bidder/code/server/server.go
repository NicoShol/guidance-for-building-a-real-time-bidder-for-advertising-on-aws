package server

import (
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
)

type Server struct {
	server fasthttp.Server
	cfg Config
}

func NewServer(
	cfg Config,
) *Server {
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		// TESTING HERE : simply log requests
		log.Info().
			Str("method", string(ctx.Method())).
			Str("path", string(ctx.Path())).
			Msg("recieved request")

		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetContentType("text/plain")
		ctx.SetBodyString("Hello from bidder!")
		// END TESTING
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