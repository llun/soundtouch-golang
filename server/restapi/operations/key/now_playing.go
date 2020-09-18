// Code generated by go-swagger; DO NOT EDIT.

package key

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// NowPlayingHandlerFunc turns a function with the right signature into a now playing handler
type NowPlayingHandlerFunc func(NowPlayingParams) middleware.Responder

// Handle executing the request and returning a response
func (fn NowPlayingHandlerFunc) Handle(params NowPlayingParams) middleware.Responder {
	return fn(params)
}

// NowPlayingHandler interface for that can handle valid now playing params
type NowPlayingHandler interface {
	Handle(NowPlayingParams) middleware.Responder
}

// NewNowPlaying creates a new http.Handler for the now playing operation
func NewNowPlaying(ctx *middleware.Context, handler NowPlayingHandler) *NowPlaying {
	return &NowPlaying{Context: ctx, Handler: handler}
}

/*NowPlaying swagger:route GET /{speakerName}/nowPlaying key device nowPlaying

This method will indicate what's playing at this moment.

*/
type NowPlaying struct {
	Context *middleware.Context
	Handler NowPlayingHandler
}

func (o *NowPlaying) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewNowPlayingParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}