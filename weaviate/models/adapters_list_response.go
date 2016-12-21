package models


// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/go-openapi/errors"
)

// AdaptersListResponse adapters list response
// swagger:model AdaptersListResponse
type AdaptersListResponse struct {

	// The list of adapters.
	Adapters []*Adapter `json:"adapters"`

	// Identifies what kind of resource this is. Value: the fixed string "weave#adaptersListResponse".
	Kind *string `json:"kind,omitempty"`
}

// Validate validates this adapters list response
func (m *AdaptersListResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAdapters(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AdaptersListResponse) validateAdapters(formats strfmt.Registry) error {

	if swag.IsZero(m.Adapters) { // not required
		return nil
	}

	for i := 0; i < len(m.Adapters); i++ {

		if swag.IsZero(m.Adapters[i]) { // not required
			continue
		}

		if m.Adapters[i] != nil {

			if err := m.Adapters[i].Validate(formats); err != nil {
				return err
			}
		}

	}

	return nil
}
