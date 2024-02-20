// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
)

// C11yVector A Vector in the Contextionary
//
// swagger:model C11yVector
type C11yVector []float32

// Validate validates this c11y vector
func (m C11yVector) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this c11y vector based on context it is used
func (m C11yVector) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}
