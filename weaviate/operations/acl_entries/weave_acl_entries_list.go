package acl_entries


// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// WeaveACLEntriesListHandlerFunc turns a function with the right signature into a weave acl entries list handler
type WeaveACLEntriesListHandlerFunc func(WeaveACLEntriesListParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn WeaveACLEntriesListHandlerFunc) Handle(params WeaveACLEntriesListParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// WeaveACLEntriesListHandler interface for that can handle valid weave acl entries list params
type WeaveACLEntriesListHandler interface {
	Handle(WeaveACLEntriesListParams, interface{}) middleware.Responder
}

// NewWeaveACLEntriesList creates a new http.Handler for the weave acl entries list operation
func NewWeaveACLEntriesList(ctx *middleware.Context, handler WeaveACLEntriesListHandler) *WeaveACLEntriesList {
	return &WeaveACLEntriesList{Context: ctx, Handler: handler}
}

/*WeaveACLEntriesList swagger:route GET /devices/{deviceId}/aclEntries aclEntries weaveAclEntriesList

Lists ACL entries.

*/
type WeaveACLEntriesList struct {
	Context *middleware.Context
	Handler WeaveACLEntriesListHandler
}

func (o *WeaveACLEntriesList) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	var Params = NewWeaveACLEntriesListParams()

	uprinc, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	var principal interface{}
	if uprinc != nil {
		principal = uprinc
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
