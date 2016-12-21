package devices


// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/weaviate/weaviate/weaviate/models"
)

/*WeaveDevicesCreateLocalAuthTokensOK Successful response

swagger:response weaveDevicesCreateLocalAuthTokensOK
*/
type WeaveDevicesCreateLocalAuthTokensOK struct {

	// In: body
	Payload *models.DevicesCreateLocalAuthTokensResponse `json:"body,omitempty"`
}

// NewWeaveDevicesCreateLocalAuthTokensOK creates WeaveDevicesCreateLocalAuthTokensOK with default headers values
func NewWeaveDevicesCreateLocalAuthTokensOK() *WeaveDevicesCreateLocalAuthTokensOK {
	return &WeaveDevicesCreateLocalAuthTokensOK{}
}

// WithPayload adds the payload to the weave devices create local auth tokens o k response
func (o *WeaveDevicesCreateLocalAuthTokensOK) WithPayload(payload *models.DevicesCreateLocalAuthTokensResponse) *WeaveDevicesCreateLocalAuthTokensOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the weave devices create local auth tokens o k response
func (o *WeaveDevicesCreateLocalAuthTokensOK) SetPayload(payload *models.DevicesCreateLocalAuthTokensResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *WeaveDevicesCreateLocalAuthTokensOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
